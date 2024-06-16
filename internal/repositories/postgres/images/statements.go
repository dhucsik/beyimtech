package images

const (
	insertStmt = `INSERT INTO images (entity_id, entity_type, image_url, content_type, filename, format)
					VALUES ($1, $2, $3, $4, $5, $6) 
					ON CONFLICT (entity_id, entity_type, content_type, format) DO NOTHING`

	selectUserStmt = `SELECT entity_id, entity_type, image_url, content_type, filename, format
						FROM images WHERE entity_id = $1 AND entity_type = 'user'`

	selectStmt = `SELECT entity_id, entity_type, image_url, content_type, filename, format
						FROM images LIMIT $1 OFFSET $2`

	getByURLStmt = `SELECT entity_id, entity_type, image_url, content_type, filename, format
						FROM images WHERE image_url = $1`

	deleteByURLStmt = `DELETE FROM images WHERE image_url = $1`

	deleteByUserIDStmt = `DELETE FROM images WHERE entity_id = $1 AND entity_type = 'user'`
)
