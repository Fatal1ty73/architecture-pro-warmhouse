package com.warmhouse.webapp.controller;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/houses")
public class HousesController {
    private final JdbcTemplate jdbcTemplate;

    public HousesController(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    @GetMapping
    public Map<String, Object> listHouses(
            @RequestParam(defaultValue = "50") int limit,
            @RequestParam(defaultValue = "0") int offset,
            @RequestParam(name = "user_id", required = false) Long userId
    ) {
        int safeLimit = Math.max(1, Math.min(limit, 200));
        int safeOffset = Math.max(0, offset);

        String where = "";
        List<Object> args = new ArrayList<>();
        if (userId != null) {
            where = " WHERE user_id = ?";
            args.add(userId);
        }

        Integer total = jdbcTemplate.queryForObject("SELECT COUNT(*) FROM houses" + where, Integer.class, args.toArray());
        args.add(safeLimit);
        args.add(safeOffset);

        List<Map<String, Object>> items = jdbcTemplate.queryForList(
                "SELECT id, user_id, address FROM houses" + where + " ORDER BY id LIMIT ? OFFSET ?",
                args.toArray()
        );

        return Map.of("items", items, "limit", safeLimit, "offset", safeOffset, "total", total == null ? 0 : total);
    }

    @GetMapping("/{houseId}")
    public ResponseEntity<?> getHouse(@PathVariable long houseId) {
        List<Map<String, Object>> rows = jdbcTemplate.queryForList(
                "SELECT id, user_id, address FROM houses WHERE id = ?",
                houseId
        );
        if (rows.isEmpty()) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "house not found"));
        }
        return ResponseEntity.ok(rows.get(0));
    }

    @PostMapping
    public ResponseEntity<?> createHouse(@RequestBody Map<String, Object> body) {
        Object userId = body.get("user_id");
        Object address = body.get("address");
        if (userId == null || address == null) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "user_id and address are required"));
        }

        Long id = jdbcTemplate.queryForObject(
                "INSERT INTO houses(user_id, address) VALUES (?, ?) RETURNING id",
                Long.class,
                Long.parseLong(userId.toString()),
                address.toString()
        );

        return ResponseEntity.status(HttpStatus.CREATED).body(Map.of(
                "id", id,
                "user_id", Long.parseLong(userId.toString()),
                "address", address
        ));
    }

    @PatchMapping("/{houseId}")
    public ResponseEntity<?> updateHouse(@PathVariable long houseId, @RequestBody Map<String, Object> body) {
        Object address = body.get("address");
        if (address == null) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "address is required"));
        }

        int updated = jdbcTemplate.update("UPDATE houses SET address = ? WHERE id = ?", address.toString(), houseId);
        if (updated == 0) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "house not found"));
        }

        return getHouse(houseId);
    }

    @DeleteMapping("/{houseId}")
    public ResponseEntity<?> deleteHouse(@PathVariable long houseId) {
        int deleted = jdbcTemplate.update("DELETE FROM houses WHERE id = ?", houseId);
        if (deleted == 0) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("code", "NOT_FOUND", "message", "house not found"));
        }
        return ResponseEntity.noContent().build();
    }
}
