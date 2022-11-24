package main

import (
	"b/ent"
	"b/ent/user"
	"context"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"
)

func main() {

	confWrite, err := pgx.ParseConfig("postgres://haha_user:secret@127.0.0.1:5433,127.0.0.1:5432/haha?sslmode=disable&target_session_attrs=primary")
	if err != nil {
		log.Fatal(err)
	}

	confRead, err := pgx.ParseConfig("postgres://haha_user:secret@127.0.0.1:5433,127.0.0.1:5432/haha?sslmode=disable&target_session_attrs=prefer-standby")
	if err != nil {
		log.Fatal(err)
	}
	wd, rd := stdlib.OpenDB(*confWrite), stdlib.OpenDB(*confRead, stdlib.OptionBeforeConnect(stdlib.RandomizeHostOrderFunc))

	//Run the auto migration tool.
	if err := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, wd))).Debug().Schema.Create(context.Background()); err != nil {
		log.Printf("failed creating schema resources: %v\n", err)
	}

	_, err = wd.Exec(`INSERT INTO "users" ("age", "name") VALUES (30, 'a8m')`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = rd.Exec(`INSERT INTO "users" ("age", "name") VALUES (30, 'a8m')`)
	if err != nil {
		log.Println("========rd=======", err)
	}
	_ = rd
	entClient := ent.NewClient(ent.Driver(&multiDriver{r: entsql.OpenDB(dialect.Postgres, rd), w: entsql.OpenDB(dialect.Postgres, wd)})).Debug()
	defer entClient.Close()

	ctx := context.Background()

	InitGorm()

	err = CreateUserGorm(ctx, db)
	if err != nil {
		log.Println("==============  CreateUserGorm  ===============", err)
	}

	user, err := QueryUserGorm(ctx, db)
	if err != nil {
		log.Println("==============  QueryUserGorm  ===============", err)
	}
	log.Println("QueryUserGorm: ", user)

	for {

		fmt.Println("try to insert")

		_, err = CreateUser(ctx, entClient)
		if err != nil {
			time.Sleep(1 * time.Second)
			fmt.Println("CreateUser error: ", err)

			continue
		}

		fmt.Println("successfully inserted!")

		fmt.Println("try to select")

		_, err = QueryUser(ctx, entClient)
		if err != nil {
			time.Sleep(1 * time.Second)

			continue
		}

		fmt.Println("successfully selected!")

		time.Sleep(5 * time.Second)
	}
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ("a8m")).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

//
//func tryToSelect(db *sql.DB) error {
//	rows, err := db.Query("select name, value from t;")
//	if err != nil {
//		fmt.Printf("failed to query data: %v\n", err)
//
//		return err
//	}
//	defer rows.Close() // nolint: errcheck
//
//	for rows.Next() {
//		var name, value string
//
//		err = rows.Scan(&name, &value)
//		if err != nil {
//			fmt.Printf("failed to scan row: %v\n", err)
//
//			continue
//		}
//
//		fmt.Printf("selected name: %s, value %s\n", name, value)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		fmt.Printf("rows error: %v\n", err)
//
//		return err
//	}
//
//	return nil
//}
//
//func tryToInsert(db *sql.DB) error {
//	const q = `insert into t (name, value) values ('Anton', '{"test2": 2}'::jsonb)`
//
//	_, err := db.Exec(q)
//	if err != nil {
//		fmt.Printf("failed to exec: %v\n", err)
//
//		return err
//	}
//
//	return nil
//}
