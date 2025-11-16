import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
    scenarios: {
        health: {
            executor: 'constant-vus',
            vus: 2,
            duration: '20s',
            exec: 'testHealth',
        },

        stats: {
            executor: 'constant-vus',
            vus: 10,
            duration: '30s',
            exec: 'testStats',
        },

        teamGet: {
            executor: 'constant-vus',
            vus: 8,
            duration: '30s',
            exec: 'testTeamGet',
        },

        teamAdd: {
            executor: 'constant-arrival-rate',
            rate: 2,
            timeUnit: '1s',
            duration: '20s',
            preAllocatedVUs: 5,
            maxVUs: 10,
            exec: 'testTeamAdd',
        },

        usersSetActive: {
            executor: 'constant-vus',
            vus: 5,
            duration: '20s',
            exec: 'testSetActive',
        },

        usersGetReview: {
            executor: 'constant-vus',
            vus: 5,
            duration: '25s',
            exec: 'testGetReview',
        },

        prCreate: {
            executor: 'constant-arrival-rate',
            rate: 3,
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 10,
            maxVUs: 20,
            exec: 'testPRCreate',
        },

        prMerge: {
            executor: 'constant-vus',
            vus: 3,
            duration: '20s',
            exec: 'testPRMerge',
        },

        prReassign: {
            executor: 'constant-vus',
            vus: 6,
            duration: '25s',
            exec: 'testPRReassign',
        },
    },
};

export function testHealth() {
    http.get('http://localhost:8080/health');
    sleep(0.1);
}

export function testStats() {
    http.get('http://localhost:8080/stats');
    sleep(0.1);
}

export function testTeamGet() {
    http.get('http://localhost:8080/team/get?team_name=backend');
    sleep(0.1);
}

export function testTeamAdd() {
    const payload = JSON.stringify({
        team_name: "team_load_" + Math.floor(Math.random() * 9999),
        members: [
            { user_id: "u" + Math.random(), username: "tester", is_active: true }
        ]
    });

    http.post('http://localhost:8080/team/add', payload, {
        headers: { "Content-Type": "application/json" }
    });

    sleep(0.2);
}

export function testSetActive() {
    const payload = JSON.stringify({
        user_id: "u1",
        is_active: true
    });

    http.post('http://localhost:8080/users/setIsActive', payload, {
        headers: { "Content-Type": "application/json" }
    });

    sleep(0.1);
}

export function testGetReview() {
    http.get('http://localhost:8080/users/getReview?user_id=u1');
    sleep(0.1);
}

export function testPRCreate() {
    const payload = JSON.stringify({
        pull_request_id: "pr_" + Math.random(),
        pull_request_name: "Load test PR",
        author_id: "u1"
    });

    http.post('http://localhost:8080/pullRequest/create', payload, {
        headers: { "Content-Type": "application/json" }
    });

    sleep(0.1);
}

export function testPRMerge() {
    const payload = JSON.stringify({ pull_request_id: "pr1" });

    http.post('http://localhost:8080/pullRequest/merge', payload, {
        headers: { "Content-Type": "application/json" }
    });

    sleep(0.2);
}

export function testPRReassign() {
    const payload = JSON.stringify({
        pull_request_id: "pr1",
        old_user_id: "u2"
    });

    http.post('http://localhost:8080/pullRequest/reassign', payload, {
        headers: { "Content-Type": "application/json" }
    });

    sleep(0.2);
}
