package com.warmhouse.webapp.controller;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/users")
public class UsersController {
    private final JdbcTemplate jdbcTemplate;

    public UsersController(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    @GetMapping
    public Map<String, Object> listUsers(
            @RequestParam(defaultValue = "50") int limit,
            @RequestParam(defaultValue = "0") int offset
    ) {
        int safeLimit = Math.max(1, Math.min(limit, 200));
        int safeOffset = Math.max(0, offset);
        Integer total = jdbcTemplate.queryForObject("SELECT COUNT(*) FROM users", Integer.class);
        List<Map<String, Object>> items = jdbcTemplate.queryForList(
                "SELECT id, name, email, login FROM users ORDER BY id LIMIT ? OFFSET ?",
                safeLimit,
                safeOffset
        );
        return Map.of("items", items, "limit", safeLimit, "offset", safeOffset, "total", total == null ? 0 : total);
    }

    @GetMapping("/{userId}")
    public ResponseEntity<?> getUser(@PathVariable long userId) {
        List<Map<String, Object>> rows = jdbcTemplate.queryForList(
                "SELECT id, name, email, login FROM users WHERE id = ?",
                userId
        );
        if (rows.isEmpty()) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "user not found"));
        }
        return ResponseEntity.ok(rows.get(0));
    }

    @PostMapping
    public ResponseEntity<?> createUser(@RequestBody Map<String, Object> body) {
        Object name = body.get("name");
        Object email = body.get("email");
        Object login = body.get("login");
        Object password = body.get("password");

        if (name == null || email == null || login == null || password == null) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "name,email,login,password are required"));
        }

        Long id = jdbcTemplate.queryForObject(
                "INSERT INTO users(name, email, login, password) VALUES (?, ?, ?, ?) RETURNING id",
                Long.class,
                name.toString(),
                email.toString(),
                login.toString(),
                password.toString()
        );

        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "id", id,
                "name", name,
                "email", email,
                "login", login
        ));
    }

    @PatchMapping("/{userId}")
    public ResponseEntity<?> updateUser(@PathVariable long userId, @RequestBody Map<String, Object> body) {
        if (body == null || body.isEmpty()) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "empty payload"));
        }

        List<String> setParts = new java.util.ArrayList<>();
        List<Object> args = new java.util.ArrayList<>();
        for (String key : List.of("name", "email", "login", "password")) {
            if (body.containsKey(key)) {
                setParts.add(key + " = ?");
                args.add(body.get(key));
            }
        }

        if (setParts.isEmpty()) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "no updatable fields"));
        }

        args.add(userId);
        int updated = jdbcTemplate.update(
                "UPDATE users SET " + String.join(", ", setParts) + " WHERE id = ?",
                args.toArray()
        );

        if (updated == 0) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "user not found"));
        }

        return getUser(userId);
    }

    @DeleteMapping("/{userId}")
    public ResponseEntity<?> deleteUser(@PathVariable long userId) {
        int deleted = jdbcTemplate.update("DELETE FROM users WHERE id = ?", userId);
        if (deleted == 0) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "user not found"));
        }
        return ResponseEntity.noContent().build();
    }
}
