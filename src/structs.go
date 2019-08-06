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
	ID                         int64       `json:"id,omitempty"`
	Name                       string      `json:"name,omitempty"`
	Priority                   string      `json:"priority,omitempty"`
	Status                     string      `json:"status,omitempty"`
	ParentTaskID               int64       `json:"parentTaskId,omitempty"`
	Description                string      `json:"description,omitempty"`
	CanViewEstTime             bool        `json:"canViewEstTime,omitempty"`
	UpdatedBy                  *AtedBy     `json:"updatedBy,omitempty"`
	CreatedBy                  *AtedBy     `json:"createdBy,omitempty"`
	DateCreated                string      `json:"dateCreated,omitempty"`
	DateChanged                string      `json:"dateChanged,omitempty"`
	HasFollowers               bool        `json:"hasFollowers,omitempty"`
	HasLoggedTime              bool        `json:"hasLoggedTime,omitempty"`
	HasReminders               bool        `json:"hasReminders,omitempty"`
	HasRemindersForUser        bool        `json:"hasRemindersForUser,omitempty"`
	HasRelativeReminders       bool        `json:"hasRelativeReminders,omitempty"`
	HasTickets                 bool        `json:"hasTickets,omitempty"`
	IsPrivate                  bool        `json:"isPrivate,omitempty"`
	InstallationID             int64       `json:"installationId,omitempty"`
	PrivacyIsInherited         bool        `json:"privacyIsInherited,omitempty"`
	LockdownID                 int         `json:"lockdownId,omitempty"`
	NumMinutesLogged           int         `json:"numMinutesLogged,omitempty"`
	NumActiveSubTasks          int         `json:"numActiveSubTasks,omitempty"`
	NumAttachments             int         `json:"numAttachments,omitempty"`
	NumComments                int         `json:"numComments,omitempty"`
	NumCommentsRead            int         `json:"numCommentsRead,omitempty"`
	NumCompletedSubTasks       int         `json:"numCompletedSubTasks,omitempty"`
	NumDependencies            int         `json:"numDependencies,omitempty"`
	NumEstMins                 int         `json:"numEstMins,omitempty"`
	NumPredecessors            int         `json:"numPredecessors,omitempty"`
	Position                   int         `json:"position,omitempty"`
	ProjectID                  int         `json:"projectId,omitempty"`
	StartDate                  interface{} `json:"startDate"`
	Tags                       []Tag       `json:"tags"`
	DueDate                    interface{} `json:"dueDate"`
	DueDateFromMilestone       bool        `json:"dueDateFromMilestone,omitempty"`
	TaskListID                 int         `json:"taskListId,omitempty"`
	Progress                   int         `json:"progress,omitempty"`
	FollowingChanges           bool        `json:"followingChanges,omitempty"`
	FollowingComments          bool        `json:"followingComments,omitempty"`
	ChangeFollowerIDS          string      `json:"changeFollowerIds,omitempty"`
	CommentFollowerIDS         string      `json:"commentFollowerIds,omitempty"`
	TeamsFollowingChanges      interface{} `json:"teamsFollowingChanges"`
	TeamsFollowingComments     interface{} `json:"teamsFollowingComments"`
	CompaniesFollowingChanges  interface{} `json:"companiesFollowingChanges"`
	CompaniesFollowingComments interface{} `json:"companiesFollowingComments"`
	Order                      int         `json:"order,omitempty"`
	CanComplete                bool        `json:"canComplete,omitempty"`
	CanEdit                    bool        `json:"canEdit,omitempty"`
	CanLogTime                 bool        `json:"canLogTime,omitempty"`
	CanAddSubtasks             bool        `json:"canAddSubtasks,omitempty"`
	Placeholder                bool        `json:"placeholder,omitempty"`
	Dlm                        int         `json:"DLM,omitempty"`
}

// AtedBy holds information abouth the person who updated and created the task
type AtedBy struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

// Tag holds information about the tag but the only tags we care about are "frontend", "backend" and "bug"
type Tag struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	ProjectID int64  `json:"projectId"`
}
