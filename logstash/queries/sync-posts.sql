SELECT l.post_log_id,
       l.operation,
       l.post_id,
       p.post_id,
       p.title,
       p.body
FROM post_logs l
         LEFT JOIN posts p
                   ON p.post_id = l.post_id
WHERE l.post_log_id > :uuid ORDER BY l.post_log_id;