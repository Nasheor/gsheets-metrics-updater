package main

// Data holds the parent data attributes returned from the API
type Data struct {
	Status    string     `json:"STATUS"`
	Tasks     []Task     `json:"tasks"`
	TaskLists []TaskList `json:"taskLists"`
	Projects  []Project  `json:"projects"`
}

// Project holds information about the project from which data is fetched from
type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	CompanyID   int64  `json:"companyId"`
	CompanyName string `json:"companyName"`
}

//TaskList holds information about task lists from which the tasks are from
type TaskList struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Task holds information about the task
type Task struct {
	ID          int64       `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	DateCreated string      `json:"dateCreated,omitempty"`
	DateChanged string      `json:"dateChanged,omitempty"`
	ProjectID   int         `json:"projectId,omitempty"`
	StartDate   interface{} `json:"startDate"`
	Tags        []Tag       `json:"tags"`
	DueDate     interface{} `json:"dueDate"`
	TaskListID  int         `json:"taskListId,omitempty"`
}

// Tag holds information about the tag but the only tags we care about are "frontend", "backend" and "bug"
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
