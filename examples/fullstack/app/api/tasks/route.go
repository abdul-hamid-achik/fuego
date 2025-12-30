package tasks

import (
	"github.com/abdul-hamid-achik/fuego/examples/fullstack/internal/tasks"
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Get returns the task list as HTML (for HTMX).
func Get(c *fuego.Context) error {
	return c.HTML(200, tasks.Default.RenderHTML())
}

// Post adds a new task.
func Post(c *fuego.Context) error {
	title := c.FormValue("title")
	if title == "" {
		return c.HTML(400, `<p class="text-red-500">Task title is required</p>`)
	}

	tasks.Default.Add(title)

	// Return the updated task list
	return c.HTML(200, tasks.Default.RenderHTML())
}

// Delete removes a task.
func Delete(c *fuego.Context) error {
	id := c.QueryInt("id", 0)
	if id == 0 {
		return c.HTML(400, `<p class="text-red-500">Task ID is required</p>`)
	}

	tasks.Default.Delete(id)

	// Return the updated task list
	return c.HTML(200, tasks.Default.RenderHTML())
}
