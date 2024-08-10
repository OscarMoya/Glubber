package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/OscarMoya/Glubber/pkg/billing"
	"github.com/OscarMoya/Glubber/pkg/model"
	"github.com/OscarMoya/Glubber/pkg/queue"
	"github.com/OscarMoya/Glubber/pkg/repository"
	"github.com/OscarMoya/Glubber/pkg/util"
	"github.com/lib/pq"
)

type RideCruder interface {
	EstimateRide(ctx context.Context, ride *model.Ride) error
	AcceptRide(ctx context.Context, ride *model.Ride) error
	DriverAccept(ctx context.Context, ride *model.Ride) error
	DriverArrived(ctx context.Context, ride *model.Ride) error
	CompleteRide(ctx context.Context, ride *model.Ride) error
	CancelRide(ctx context.Context, ride *model.Ride) error
	RideError(ctx context.Context, ride *model.Ride) error

	CreateRide(ctx context.Context, ride *model.Ride) error
	ListRides(ctx context.Context) ([]model.Ride, error)
	GetRide(ctx context.Context, id int) (*model.Ride, error)
	UpdateRide(ctx context.Context, ride *model.Ride) error
	DeleteRide(ctx context.Context, id int) error
}

type RideServiceOpts struct {
	Repository  repository.Repository
	Producer    queue.Producer
	Biller      billing.Biller
	Table       string
	DriverTopic string
	DriverKey   string
}

type RideService struct {
	RideServiceOpts
	outboxTable   string
	notifyChannel string
}

// NewRideService creates a new RideDatabase
func NewRideService(ctx context.Context, opts RideServiceOpts) (*RideService, error) {

	svc := &RideService{

		RideServiceOpts: opts,
		outboxTable:     opts.Table + "_outbox",
		notifyChannel:   opts.Table + "_events",
	}

	if err := svc.createTables(ctx); err != nil {
		return nil, err
	}

	return svc, nil
}

func (svc *RideService) processNotificationsWorker(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case notification := <-svc.Repository.Notifications(ctx):
			log.Printf("Notification: %v\n", notification)
			if notification == nil {
				continue
			}
			err := svc.processNotification(ctx, notification)
			if err != nil {
				log.Printf("Error processing notification: %v\n", err)
			}
		}
	}
}

func (svc *RideService) processNotification(ctx context.Context, notification *pq.Notification) error {
	id, err := strconv.Atoi(notification.Extra)
	if err != nil {
		return err
	}
	outbox, err := svc.getOutbox(ctx, id)
	if err != nil {
		return err
	}
	if outbox == nil {
		return fmt.Errorf("outbox not found")
	}
	switch outbox.Status {
	case model.RideStatusPassengerAccepted:
		outboxBytes, err := json.Marshal(outbox)
		if err != nil {
			return err
		}
		err = svc.Producer.SendMessage(ctx, svc.DriverTopic, svc.DriverKey, outboxBytes)
		if err != nil {
			return err
		}

		return svc.deleteOutbox(ctx, id)
	default:
		return svc.deleteOutbox(ctx, id)
	}
}

func (svc *RideService) createTables(ctx context.Context) error {
	query, err := util.BuildSQLCreateTableQuery(svc.Table, model.Ride{})
	if err != nil {
		return err
	}
	err = svc.Repository.CreateTable(ctx, query)
	if err != nil {
		return err
	}

	query, err = util.BuildSQLCreateTableQuery(svc.outboxTable, model.RideOutbox{})
	if err != nil {
		return err
	}
	err = svc.Repository.CreateTable(ctx, query)
	return err
}

func (svc *RideService) Close() {

}

func (svc *RideService) Start(ctx context.Context) error {
	err := svc.Repository.Listen(ctx, svc.notifyChannel)
	if err != nil {
		return err
	}
	go svc.processNotificationsWorker(ctx)
	return nil
}

func (svc *RideService) EstimateRide(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusPending
	svc.Biller.EstimateRide(ride)
	return svc.CreateRide(ctx, ride)
}

func (svc *RideService) AcceptRide(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusPassengerAccepted
	// After this update, the notification will be sent to the driver
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) DriverAccept(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusDriverAccepted
	// After this update, we need to notify the passengers that the driver has accepted the ride
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) DriverArrived(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusPickingUp
	// After this update, we need to notify the passengers that the driver has arrived
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) CompleteRide(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusCompleted
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) CancelRide(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusPassengerCancelled
	// After this update, we need to notify the driver that the passenger has cancelled the ride
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) RideError(ctx context.Context, ride *model.Ride) error {
	ride.Status = model.RideStatusErrored
	// After this update, we need to notify the passenger and the Driver that there was an error
	return svc.UpdateRide(ctx, ride)
}

func (svc *RideService) CreateRide(ctx context.Context, ride *model.Ride) error {
	fields, placeholder, args, _ := util.BuildSQLInsertQuery(ride, 1)
	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id;`, svc.Table, fields, placeholder)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&ride.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	outbox := model.NewRideOutbox(ride)
	fields, placeholder, args, _ = util.BuildSQLInsertQuery(outbox, 1)
	query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id;`, svc.outboxTable, fields, placeholder)
	err = tx.QueryRow(ctx, query, args...).Scan(&outbox.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	query = fmt.Sprintf(`NOTIFY %s, '%d';`, svc.notifyChannel, outbox.ID)
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
	}

	err = tx.Commit(ctx)
	return err
}

func (svc *RideService) ListRides(ctx context.Context) ([]model.Ride, error) {

	fields, _, _, _ := util.BuildSQLSelectQuery(&model.Ride{}, 1)
	query := fmt.Sprintf(`SELECT %s FROM %s;`, fields, svc.Table)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rides := make([]model.Ride, 0)
	for rows.Next() {
		ride := model.Ride{}
		err = ride.Scan(rows)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		rides = append(rides, ride)
	}
	tx.Commit(ctx)

	return rides, nil
}

func (svc *RideService) GetRide(ctx context.Context, id int) (*model.Ride, error) {
	fields, _, _, _ := util.BuildSQLSelectQuery(&model.Ride{}, 1)
	// fields = "id, " + fields
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1;`, fields, svc.Table)
	fmt.Println(query)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, query, id)
	ride := &model.Ride{}
	err = ride.Scan(row)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Commit(ctx)
	return ride, nil
}

func (svc *RideService) UpdateRide(ctx context.Context, ride *model.Ride) error {
	setStmt, args, _ := util.BuildSQLUpdateQuery(ride, 1)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d;`, svc.Table, setStmt, ride.ID)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	outbox := model.NewRideOutbox(ride)
	fields, placeholder, args, _ := util.BuildSQLInsertQuery(outbox, 1)
	query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id;`, svc.outboxTable, fields, placeholder)
	err = tx.QueryRow(ctx, query, args...).Scan(&outbox.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	query = fmt.Sprintf(`NOTIFY %s, '%d';`, svc.notifyChannel, outbox.ID)
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (svc *RideService) DeleteRide(ctx context.Context, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = %d;`, svc.Table, id)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	outbox := &model.RideOutbox{RideID: id, Status: model.RideStatusDeleted}
	fields, placeholder, args, _ := util.BuildSQLInsertQuery(outbox, 1)
	query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s) RETURNING id;`, svc.outboxTable, fields, placeholder)
	err = tx.QueryRow(ctx, query, args...).Scan(&outbox.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	query = fmt.Sprintf(`NOTIFY %s, '%d';`, svc.notifyChannel, outbox.ID)
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (svc *RideService) DeleteAllRides(ctx context.Context) error {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, svc.Table)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	query = fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, svc.outboxTable)
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	log.Println("Deleted all rides ERR", err)
	return err

}

func (svc *RideService) getOutbox(ctx context.Context, id int) (*model.RideOutbox, error) {
	fields, _, _, _ := util.BuildSQLSelectQuery(&model.RideOutbox{}, 1)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1;`, fields, svc.outboxTable)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, query, id)
	outbox := &model.RideOutbox{}
	err = outbox.Scan(row)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	tx.Commit(ctx)
	return outbox, nil
}

func (svc *RideService) deleteOutbox(ctx context.Context, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = %d;`, svc.outboxTable, id)
	tx, err := svc.Repository.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
