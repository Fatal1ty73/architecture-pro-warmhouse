package com.warmhouse.webapp.controller;

import com.warmhouse.webapp.auth.AuthService;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/auth")
public class AuthController {
    private final AuthService authService;

    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    @GetMapping("/login")
    public ResponseEntity<?> login(
            @RequestParam(name = "username", required = true) String username,
            @RequestParam(name = "password", required = true) String password
    ) {
        if (username == null || password == null) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "username and password are required"));
        }

        String token = authService.login(username, password);
        if (token == null) {
            return ResponseEntity.badRequest().body(Map.of("code", "BAD_REQUEST", "message", "invalid username/password"));
        }

        HttpHeaders headers = new HttpHeaders();
        headers.add("X-Rate-Limit", "1000");
        headers.add("X-Expires-After", DateTimeFormatter.ISO_OFFSET_DATE_TIME
                .format(Instant.ofEpochSecond(authService.expiresAt(token)).atOffset(ZoneOffset.UTC)));
        return new ResponseEntity<>(token, headers, HttpStatus.OK);
    }

    @GetMapping("/logout")
    public ResponseEntity<?> logout(@RequestHeader(name = "X-Auth-Token", required = false) String token) {
        authService.logout(token);
        return ResponseEntity.ok(Map.of("message", "logged out"));
    }
}
