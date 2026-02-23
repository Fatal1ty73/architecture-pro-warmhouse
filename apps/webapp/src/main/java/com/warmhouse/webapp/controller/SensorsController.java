package com.warmhouse.webapp.controller;

import com.warmhouse.webapp.service.ProxyClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.LinkedHashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/sensors")
public class SensorsController {
    private final ProxyClient proxyClient;
    private final String sensorManagerUrl;

    public SensorsController(ProxyClient proxyClient,
                             @Value("${clients.sensor-manager-url}") String sensorManagerUrl) {
        this.proxyClient = proxyClient;
        this.sensorManagerUrl = sensorManagerUrl;
    }

    @GetMapping
    public ResponseEntity<String> listSensors(
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
        return proxyClient.forward(sensorManagerUrl + "/api/v1/sensors", HttpMethod.GET, query, null);
    }

    @GetMapping("/{sensorId}")
    public ResponseEntity<String> getSensor(@PathVariable String sensorId) {
        return proxyClient.forward(sensorManagerUrl + "/api/v1/sensors/" + sensorId, HttpMethod.GET, null, null);
    }

    @PostMapping
    public ResponseEntity<String> createSensor(@RequestBody String body) {
        return proxyClient.forward(sensorManagerUrl + "/api/v1/sensors", HttpMethod.POST, null, body);
    }

    @PatchMapping("/{sensorId}")
    public ResponseEntity<String> updateSensor(@PathVariable String sensorId, @RequestBody String body) {
        return proxyClient.forward(sensorManagerUrl + "/api/v1/sensors/" + sensorId, HttpMethod.PATCH, null, body);
    }

    @DeleteMapping("/{sensorId}")
    public ResponseEntity<String> deleteSensor(@PathVariable String sensorId) {
        return proxyClient.forward(sensorManagerUrl + "/api/v1/sensors/" + sensorId, HttpMethod.DELETE, null, null);
    }
}
