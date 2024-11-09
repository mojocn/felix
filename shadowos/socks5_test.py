import socket
import socks
import json
import requests


def test_dns_query():
    # Proxy details
    proxy_host = '127.0.0.1'
    proxy_port = 2080

    # DNS server details
    dns_server = '114.114.114.114'  # Example: Google DNS
    dns_port = 53

    # Create a UDP socket with SOCKS5 proxy
    sock = socks.socksocket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.set_proxy(socks.SOCKS5, proxy_host, proxy_port)
    # DNS query message
    # This is a simple DNS query for the domain "example.com"
    # Transaction ID: 0x1234
    # Flags: Standard query
    # Questions: 1
    # Answer RRs: 0
    # Authority RRs: 0
    # Additional RRs: 0
    # Query: example.com, Type: A, Class: IN
    message = b'\x12\x34'  # Transaction ID
    message += b'\x01\x00'  # Flags
    message += b'\x00\x01'  # Questions
    message += b'\x00\x00'  # Answer RRs
    message += b'\x00\x00'  # Authority RRs
    message += b'\x00\x00'  # Additional RRs
    message += b'\x07example\x03com\x00'  # Query: example.com
    message += b'\x00\x01'  # Type: A
    message += b'\x00\x01'  # Class: IN

    # Send the DNS query
    sock.sendto(message, (dns_server, dns_port))

    # Receive a response
    response, addr = sock.recvfrom(4096)
    print(f'Received response from {addr}: {response}')

    # Close the socket
    sock.close()


def test_http(size=1 * 1024):
    # SOCKS5 proxy configuration
    proxies = {
        'http': 'socks5h://localhost:2080',
        'https': 'socks5h://localhost:2080'
    }

    # Generate a 2MB JSON body
    data = {'key': 'a' * size}  # Adjusting for JSON formatting characters

    # Convert the data to JSON format
    json_data = json.dumps(data)

    # URL to send the POST request to
    url = 'http://httpbin.org/post'

    # Send the POST request
    response = requests.post(url, data=json_data, headers={'Content-Type': 'application/json'}, proxies=proxies)

    # Print the response
    print(f'Status Code: {response.status_code}')
    print(f'Response Body: {response.text}')



# test_via_http()

test_http(5)
# test_dns_query()
# test_dns_query()
# test_http()
