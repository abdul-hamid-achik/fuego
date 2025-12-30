package toggle

import (
	"github.com/abdul-hamid-achik/fuego/examples/fullstack/internal/tasks"
	"github.com/abdul-hamid-achik/fuego/pkg/fuego"
)

// Post toggles the completion status of a task.
func Post(c *fuego.Context) error {
	id := c.QueryInt("id", 0)
	if id == 0 {
		return c.HTML(400, `<p class="text-red-500">Task ID is required</p>`)
	}

	tasks.Default.Toggle(id)

	// Return the updated task list
	return c.HTML(200, tasks.Default.RenderHTML())
}
