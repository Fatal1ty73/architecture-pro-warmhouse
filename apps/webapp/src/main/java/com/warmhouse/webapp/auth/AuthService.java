package com.warmhouse.webapp.auth;

import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class AuthService {
    private static final long TOKEN_TTL_SECONDS = 60 * 60 * 8;

    private final JdbcTemplate jdbcTemplate;
    private final ConcurrentHashMap<String, Session> sessions = new ConcurrentHashMap<>();

    public AuthService(JdbcTemplate jdbcTemplate) {
        this.jdbcTemplate = jdbcTemplate;
    }

    public String login(String username, String password) {
        Integer count = jdbcTemplate.queryForObject(
                "SELECT COUNT(*) FROM users WHERE login = ? AND password = ?",
                Integer.class,
                username,
                password
        );
        if (count == null || count == 0) {
            return null;
        }

        String token = UUID.randomUUID().toString();
        long expiresAt = Instant.now().plusSeconds(TOKEN_TTL_SECONDS).getEpochSecond();
        sessions.put(token, new Session(username, expiresAt));
        return token;
    }

    public boolean isValid(String token) {
        if (token == null || token.isBlank()) {
            return false;
        }
        Session session = sessions.get(token);
        if (session == null) {
            return false;
        }
        if (session.expiresAtEpochSeconds < Instant.now().getEpochSecond()) {
            sessions.remove(token);
            return false;
        }
        return true;
    }

    public void logout(String token) {
        if (token != null) {
            sessions.remove(token);
        }
    }

    public long expiresAt(String token) {
        Session session = sessions.get(token);
        return session == null ? 0L : session.expiresAtEpochSeconds;
    }

    private record Session(String username, long expiresAtEpochSeconds) {}
}
