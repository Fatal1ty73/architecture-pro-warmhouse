package com.warmhouse.webapp.controller;

import com.warmhouse.webapp.service.ProxyClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.LinkedHashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/devices")
public class DevicesController {
    private final ProxyClient proxyClient;
    private final String deviceHandleUrl;

    public DevicesController(ProxyClient proxyClient,
                             @Value("${clients.device-handle-url}") String deviceHandleUrl) {
        this.proxyClient = proxyClient;
        this.deviceHandleUrl = deviceHandleUrl;
    }

    @GetMapping
    public ResponseEntity<String> listDevices(
            @RequestParam(required = false) String limit,
            @RequestParam(required = false) String offset,
            @RequestParam(name = "house_id", required = false) String houseId,
            @RequestParam(name = "type_id", required = false) String typeId,
            @RequestParam(required = false) String status
    ) {
        Map<String, String> query = new LinkedHashMap<>();
        query.put("limit", limit);
        query.put("offset", offset);
        query.put("house_id", houseId);
        query.put("type_id", typeId);
        query.put("status", status);
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices", HttpMethod.GET, query, null);
    }

    @GetMapping("/{deviceId}")
    public ResponseEntity<String> getDevice(@PathVariable String deviceId) {
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices/" + deviceId, HttpMethod.GET, null, null);
    }

    @PostMapping
    public ResponseEntity<String> createDevice(@RequestBody String body) {
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices", HttpMethod.POST, null, body);
    }

    @PatchMapping("/{deviceId}")
    public ResponseEntity<String> updateDevice(@PathVariable String deviceId, @RequestBody String body) {
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices/" + deviceId, HttpMethod.PATCH, null, body);
    }

    @DeleteMapping("/{deviceId}")
    public ResponseEntity<String> deleteDevice(@PathVariable String deviceId) {
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices/" + deviceId, HttpMethod.DELETE, null, null);
    }

    @PostMapping("/{deviceId}/execute")
    public ResponseEntity<String> executeDevice(@PathVariable String deviceId, @RequestBody(required = false) String body) {
        return proxyClient.forward(deviceHandleUrl + "/api/v1/devices/" + deviceId + "/execute", HttpMethod.POST, null, body);
    }
}
