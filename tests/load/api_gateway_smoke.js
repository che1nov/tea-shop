import http from "k6/http";
import { check, sleep } from "k6";

const baseURL = __ENV.BASE_URL || "http://localhost:8080";

export const options = {
  vus: 5,
  duration: "30s",
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<800"],
    checks: ["rate>0.99"],
  },
};

export default function () {
  const res = http.get(`${baseURL}/swagger/index.html`, {
    tags: { endpoint: "swagger" },
    timeout: "10s",
  });

  check(res, {
    "swagger status is 200": (r) => r.status === 200,
  });

  sleep(0.1);
}
