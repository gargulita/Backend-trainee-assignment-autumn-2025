import http from "k6/http";
import { sleep } from "k6";

export let options = {
    vus: 20,          
    duration: "20s", 
};

export default function () {
    const url = "http://localhost:8080/team/deactivate";
    const payload = JSON.stringify({
        team_name: "backend"
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
        },
    };

    http.post(url, payload, params);

    sleep(0.1);
}
