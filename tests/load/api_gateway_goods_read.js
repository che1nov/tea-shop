import http from "k6/http";
import { check, sleep } from "k6";

const baseURL = __ENV.BASE_URL || "http://localhost:8080";
const duration = __ENV.DURATION || "3m";
const vus = Number(__ENV.VUS || 50);
const thinkTime = Number(__ENV.THINK_TIME_SEC || 0.2);

export const options = {
  scenarios: {
    goods_read_load: {
      executor: "constant-vus",
      vus,
      duration,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500", "p(99)<1000"],
    checks: ["rate>0.99"],
  },
};

export default function () {
  const limit = 10 + Math.floor(Math.random() * 20);
  const offset = Math.floor(Math.random() * 200);
  const url = `${baseURL}/api/v1/goods?limit=${limit}&offset=${offset}`;

  const res = http.get(url, {
    tags: { endpoint: "list_goods" },
    timeout: "10s",
  });

  check(res, {
    "goods list status is 200": (r) => r.status === 200,
    "goods list response is json": (r) =>
      (r.headers["Content-Type"] || "").includes("application/json"),
  });

  sleep(thinkTime);
}
