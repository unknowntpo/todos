/*
 * TODOS API
 *
 * Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.
 *
 * API version: 0.1.9
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type CreateTaskRequest struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Done bool `json:"done"`
}