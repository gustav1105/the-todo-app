package service

import (
    "context"
    "database/sql"
    "fmt"

    pb "the-todo-app/proto"
    "google.golang.org/protobuf/types/known/timestamppb"
)

// TodoService implements pb.TodoServiceServer
type TodoService struct {
    pb.UnimplementedTodoServiceServer
    DB *sql.DB
}

// NewTodoService creates a new instance of TodoService
func NewTodoService(db *sql.DB) *TodoService {
    return &TodoService{DB: db}
}

// AddTask handles adding a new task to the database
func (s *TodoService) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
    task := req.GetTask()

    // Prepare SQL to insert a new task
    stmt, err := s.DB.Prepare("INSERT INTO tasks(name, note, status, due_on) VALUES (?, ?, ?, ?)")
    if err != nil {
        return nil, fmt.Errorf("failed to prepare SQL statement: %v", err)
    }
    defer stmt.Close()

    // Execute the query
    res, err := stmt.Exec(task.Name, task.Note, task.Status, task.DueOn.AsTime())
    if err != nil {
        return nil, fmt.Errorf("failed to execute SQL statement: %v", err)
    }

    // Retrieve the last insert ID (task ID)
    id, err := res.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve last insert ID: %v", err)
    }

    task.Id = int32(id)

    return &pb.AddTaskResponse{
        Task: task,
    }, nil
}

// GetTask retrieves a task by its ID
func (s *TodoService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
    task := &pb.Task{}
    id := req.GetId()

    // Query for the task by ID
    err := s.DB.QueryRow("SELECT id, name, note, status, due_on FROM tasks WHERE id = ?", id).
        Scan(&task.Id, &task.Name, &task.Note, &task.Status, &task.DueOn)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("task with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to query task: %v", err)
    }

    return &pb.GetTaskResponse{
        Task: task,
    }, nil
}

// GetTasks retrieves all tasks from the database
func (s *TodoService) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
    rows, err := s.DB.Query("SELECT id, name, note, status, due_on FROM tasks")
    if err != nil {
        return nil, fmt.Errorf("failed to query tasks: %v", err)
    }
    defer rows.Close()

    var tasks []*pb.Task

    // Iterate through each row
    for rows.Next() {
        task := &pb.Task{}
        var dueOn sql.NullTime
        if err := rows.Scan(&task.Id, &task.Name, &task.Note, &task.Status, &dueOn); err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }

        if dueOn.Valid {
            task.DueOn = timestamppb.New(dueOn.Time)
        }

        tasks = append(tasks, task)
    }

    return &pb.GetTasksResponse{
        Tasks: tasks,
    }, nil
}

// CompleteTask marks a task as completed
func (s *TodoService) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
    id := req.GetId()

    // Update task status to COMPLETED
    _, err := s.DB.Exec("UPDATE tasks SET status = ? WHERE id = ?", pb.TaskStatus_COMPLETED, id)
    if err != nil {
        return nil, fmt.Errorf("failed to update task status: %v", err)
    }

    // Query the updated task
    task := &pb.Task{}
    err = s.DB.QueryRow("SELECT id, name, note, status, due_on FROM tasks WHERE id = ?", id).
        Scan(&task.Id, &task.Name, &task.Note, &task.Status, &task.DueOn)
    if err != nil {
        return nil, fmt.Errorf("failed to query updated task: %v", err)
    }

    return &pb.CompleteTaskResponse{
        Task: task,
    }, nil
}

// CancelTask marks a task as canceled
func (s *TodoService) CancelTask(ctx context.Context, req *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
    id := req.GetId()

    // Update task status to CANCELLED
    _, err := s.DB.Exec("UPDATE tasks SET status = ? WHERE id = ?", pb.TaskStatus_CANCELLED, id)
    if err != nil {
        return nil, fmt.Errorf("failed to update task status: %v", err)
    }

    // Query the updated task
    task := &pb.Task{}
    err = s.DB.QueryRow("SELECT id, name, note, status, due_on FROM tasks WHERE id = ?", id).
        Scan(&task.Id, &task.Name, &task.Note, &task.Status, &task.DueOn)
    if err != nil {
        return nil, fmt.Errorf("failed to query updated task: %v", err)
    }

    return &pb.CancelTaskResponse{
        Task: task,
    }, nil
}

