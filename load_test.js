import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

// Read the client cert and key for mTLS
const clientCert = open('./certs/client.crt');
const clientKey = open('./certs/client.key');

export const options = {
    tlsAuth: [
        {
            domains: ['service1.local', 'service2.local', 'localhost'],
            cert: clientCert,
            key: clientKey,
        },
    ],
    insecureSkipTLSVerify: true, // we are using a self-signed CA for testing
    vus: 50,
    duration: '10s',
    hosts: {
        'service1.local': '127.0.0.1',
        'service2.local': '127.0.0.1',
    },
};

export default function () {
    // We use k6's host resolution feature to map 'service1.local' to 127.0.0.1
    // This allows us to use the actual SNI in the request URL.

    const res1 = http.get('https://service1.local:8443/', {
        headers: { 'Host': 'service1.local' }
    });

    // For HTTP, host header isn't always enough to change TLS SNI in some clients,
    // but we will verify if it works or if a direct SNI test is needed.

    check(res1, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(0.1);
}
