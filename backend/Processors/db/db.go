package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./LHS.db")
	if err != nil {
		return err
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	createTables := `
	CREATE TABLE IF NOT EXISTS stacks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS stack_services (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id INTEGER,
		container_id TEXT,
		name TEXT,
		image TEXT,
		build_path TEXT,
		build_dockerfile TEXT,
		ports TEXT,
		env TEXT,
		volumes TEXT,
		FOREIGN KEY(stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);`
	_, err = DB.Exec(createTables)
	if err != nil {
		return err
	}
	return nil
}

//STACK CRUD

func CreateStack(name string) (int64, error) {
	res, err := DB.Exec("INSERT INTO stacks (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetStackById(id int64) (*Stack, error) {
	var stack Stack
	err := DB.QueryRow("SELECT id, name, created_at FROM stacks WHERE id = ?", id).Scan(&stack.ID, &stack.Name, &stack.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

func GetAllStacks() ([]Stack, error) {
	rows, err := DB.Query("SELECT id, name, created_at FROM stacks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []Stack
	for rows.Next() {
		var s Stack
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		stacks = append(stacks, s)
	}
	return stacks, nil
}

func DeleteStack(id int64) error {
	_, err := DB.Exec("DELETE FROM stacks WHERE id = ?", id)
	return err
}

func UpdateStack(id int64, newName string) error {
	_, err := DB.Exec("UPDATE stacks SET name = ? WHERE id = ?", newName, id)
	return err
}

func DeleteAllStacks() error {
	_, err := DB.Exec("DELETE FROM stacks")
	return err
}

//STACK SERVICE CRUD

func CreateStackService(svc StackService) error {
	_, err := DB.Exec(`
		INSERT INTO stack_services 
		(stack_id, container_id, name, image, build_path, build_dockerfile, ports, env, volumes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		svc.StackID, svc.ContainerID, svc.Name, svc.Image, svc.BuildPath, svc.BuildDockerfile,
		svc.Ports, svc.Env, svc.Volumes,
	)
	return err
}

func GetServicesByStackID(stackID int64) ([]StackService, error) {
	rows, err := DB.Query("SELECT id, stack_id, container_id, name, image, build_path, build_dockerfile, ports, env, volumes FROM stack_services WHERE stack_id = ?", stackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []StackService
	for rows.Next() {
		var svc StackService
		if err := rows.Scan(
			&svc.ID, &svc.StackID, &svc.ContainerID, &svc.Name, &svc.Image,
			&svc.BuildPath, &svc.BuildDockerfile, &svc.Ports, &svc.Env, &svc.Volumes,
		); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, rows.Err()
}

func DeleteServicesByStackID(stackID int64) error {
	_, err := DB.Exec("DELETE FROM stack_services WHERE stack_id = ?", stackID)
	return err
}

func DeleteServiceByID(serviceID int64) error {
	_, err := DB.Exec("DELETE FROM stack_services WHERE id = ?", serviceID)
	return err
}

func UpdateService(svc StackService) error {
	_, err := DB.Exec(`
		UPDATE stack_services 
		SET 
			stack_id = ?, 
			container_id = ?, 
			name = ?, 
			image = ?, 
			build_path = ?, 
			build_dockerfile = ?, 
			ports = ?, 
			env = ?, 
			volumes = ?
		WHERE id = ?`,
		svc.StackID, svc.ContainerID, svc.Name, svc.Image,
		svc.BuildPath, svc.BuildDockerfile, svc.Ports, svc.Env, svc.Volumes,
		svc.ID,
	)
	return err
}

func GetAllServices() ([]StackService, error) {
	rows, err := DB.Query("SELECT id, stack_id, container_id, name, image, build_path, build_dockerfile, ports, env, volumes FROM stack_services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []StackService
	for rows.Next() {
		var svc StackService
		if err := rows.Scan(
			&svc.ID, &svc.StackID, &svc.ContainerID, &svc.Name, &svc.Image,
			&svc.BuildPath, &svc.BuildDockerfile, &svc.Ports, &svc.Env, &svc.Volumes,
		); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, rows.Err()
}
