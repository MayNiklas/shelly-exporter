# -*- coding: utf-8 -*-
import json
import time

import requests as requests
from prometheus_client import start_http_server, Summary, Gauge

metrics_port = 8001
shelly_ip = '192.168.15.2'
request_interval = 10

# metrics
REQUEST_TIME = Summary('request_processing_seconds', 'Time spent processing request')
shelly_power = Gauge('shelly_power', 'current power consumption')
shelly_power_total = Gauge('shelly_power_total', 'total power consumption')


def shelly_request_status(ip):
    response = requests.get(f"http://{ip}/status").text
    response_text = json.loads(response)
    return response_text


def shelly_request_settings(ip):
    response = requests.get(f"http://{ip}/settings").text
    response_text = json.loads(response)
    return response_text


# Decorate function with metric.
@REQUEST_TIME.time()
def process_request():
    """Request shelly metrics"""
    shelly_status = shelly_request_status(shelly_ip)
    shelly_settings = shelly_request_settings(shelly_ip)

    shelly_power.set(shelly_status['meters'][0]['power'])
    shelly_power_total.set(shelly_status['meters'][0]['total'])


def main():
    # Start up the server to expose the metrics.
    start_http_server(metrics_port)

    # Generate some requests.
    while True:
        # update metrics every { request_interval } seconds
        process_request()
        time.sleep(request_interval)


if __name__ == '__main__':
    main()
