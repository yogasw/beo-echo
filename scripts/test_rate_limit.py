#!/usr/bin/env python3
"""
API Testing Script for Beo Echo Rate Limiting

This script tests rate limiting functionality with detailed logging.
"""

import requests
import time
from datetime import datetime
from concurrent.futures import ThreadPoolExecutor, as_completed

# Configuration
BASE_URL = "http://localhost:3600"
AUTH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiODY3N2RmYzctZGU1Yy00ODk2LThmNzAtYzE3OTU5ODEyZGYyIiwiZW1haWwiOiJhZG1pbkBhZG1pbi5jb20iLCJuYW1lIjoiQWRtaW4iLCJleHAiOjE3NTA2MTQ1NTYsIm5iZiI6MTc1MDUyODE1NiwiaWF0IjoxNzUwNTI4MTU2fQ.y5BcXuoFiPiJsafALkbC_hPzaAVSGxXLvPxI-pJeGiQ"

# Test Configuration for /demo endpoint (general rate limit)
DEMO_SUKSES_EXPECTED = 200  # Expected successful requests for /demo
DEMO_GAGAL_EXPECTED = 50    # Expected failed/rate limited requests for /demo
DEMO_TOTAL_REQUESTS = DEMO_SUKSES_EXPECTED + DEMO_GAGAL_EXPECTED  # Total requests to send

# Test Configuration for /api endpoint (API rate limit)
API_SUKSES_EXPECTED = 60    # Expected successful requests for /api
API_GAGAL_EXPECTED = 50     # Expected failed/rate limited requests for /api
API_TOTAL_REQUESTS = API_SUKSES_EXPECTED + API_GAGAL_EXPECTED   # Total requests to send

def get_headers():
    """Get common headers for API requests"""
    return {
        'Accept': 'application/json, text/plain, */*',
        'Authorization': f'Bearer {AUTH_TOKEN}',
        'Content-Type': 'application/json',
        'User-Agent': 'Python-Test-Script'
    }

def wait_for_next_minute():
    """Wait until the next minute starts at second 1"""
    now = datetime.now()
    current_second = now.second
    
    if current_second > 1:
        # Calculate seconds to wait until next minute + 1 second
        wait_seconds = 60 - current_second + 1
        print(f"⏰ Menunggu {wait_seconds} detik untuk mulai di detik ke-1 menit baru...")
        time.sleep(wait_seconds)
    
    # Get exact start time
    start_time = datetime.now()
    print(f"✅ Mulai test pada {start_time.strftime('%H:%M:%S')} (detik ke-{start_time.second})")
    return start_time

def test_single_request(url, request_num, use_auth=False):
    """Make a single request and return result"""
    try:
        headers = {'User-Agent': 'Python-Test-Script'}
        if use_auth:
            headers.update(get_headers())
            
        response = requests.get(url, headers=headers, timeout=5)
        timestamp = datetime.now().strftime('%H:%M:%S.%f')[:-3]
        
        return {
            'request_num': request_num,
            'status_code': response.status_code,
            'timestamp': timestamp,
            'success': response.status_code != 429,
            'rate_limited': response.status_code == 429
        }
    except Exception as e:
        return {
            'request_num': request_num,
            'error': str(e),
            'timestamp': datetime.now().strftime('%H:%M:%S.%f')[:-3],
            'success': False,
            'rate_limited': False
        }

def test_rate_limiting_sequential(endpoint, num_requests, use_auth=False):
    """Test rate limiting with sequential requests and detailed logging"""
    expected_limit = API_SUKSES_EXPECTED if '/api' in endpoint else DEMO_SUKSES_EXPECTED
    expected_failures = num_requests - expected_limit
    
    print(f"🧪 Testing {endpoint}")
    print(f"   Expected: {expected_limit} sukses, {expected_failures} gagal")
    print(f"   Total requests: {num_requests}")
    
    # Wait for clean minute start
    start_time = wait_for_next_minute()
    
    url = f"{BASE_URL}{endpoint}"
    results = []
    success_count = 0
    rate_limited_count = 0
    error_count = 0
    
    # Send requests sequentially with logging
    for i in range(num_requests):
        request_num = i + 1
        result = test_single_request(url, request_num, use_auth)
        results.append(result)
        
        # Log each request
        if 'error' in result:
            error_count += 1
            print(f"   ❌ Request {request_num:3d} [{result['timestamp']}]: ERROR - {result['error']}")
        elif result['rate_limited']:
            rate_limited_count += 1
            print(f"   🛑 Request {request_num:3d} [{result['timestamp']}]: RATE LIMITED (429)")
        else:
            success_count += 1
            print(f"   ✅ Request {request_num:3d} [{result['timestamp']}]: SUCCESS ({result['status_code']})")
        
        # Small delay to prevent overwhelming the server
        time.sleep(0.01)
    
    # Print final results
    end_time = datetime.now()
    total_time = (end_time - start_time).total_seconds()
    
    print(f"\n📊 HASIL TEST {endpoint}")
    print(f"   Total requests: {num_requests}")
    print(f"   Waktu total: {total_time:.2f} detik")
    print(f"   ✅ SUKSES: {success_count}")
    print(f"   🛑 RATE LIMITED: {rate_limited_count}")
    if error_count > 0:
        print(f"   ❌ ERROR: {error_count}")
    
    # Final verdict
    if success_count == expected_limit and rate_limited_count >= 1:
        print(f"   🎯 HASIL: BENAR! {expected_limit} sukses, {rate_limited_count} gagal")
    elif success_count == expected_limit:
        print(f"   ⚠️  HASIL: {expected_limit} sukses tapi tidak ada yang rate limited")
    elif rate_limited_count >= 1:
        print(f"   ⚠️  HASIL: Hanya {success_count} sukses (harusnya {expected_limit}), {rate_limited_count} gagal")
    else:
        print(f"   ❌ HASIL: SALAH! {success_count} sukses, {rate_limited_count} gagal (harusnya {expected_limit} sukses, 1+ gagal)")
    
    print("-" * 60)
    return results

def test_server_connectivity():
    """Test if the server is running and accessible"""
    print("🔌 Testing server connectivity...")
    
    try:
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        if response.status_code == 200:
            print(f"   ✅ Server is running at {BASE_URL}")
            return True
        else:
            print(f"   ⚠️  Server responded with status {response.status_code}")
            return True  # Server is running but maybe no health endpoint
    except requests.exceptions.ConnectionError:
        print(f"   ❌ Cannot connect to server at {BASE_URL}")
        print("   Make sure the server is running on port 3600")
        return False
    except Exception as e:
        print(f"   ❌ Unexpected error: {e}")
        return False

def main():
    """Main function to run all tests"""
    print("🚀 Testing Rate Limiting Beo Echo")
    print("=" * 50)
    
    # Test server connectivity first
    if not test_server_connectivity():
        print("\n❌ Server tidak dapat diakses")
        return
    
    print()
    
    # Test 1: Rate limiting on /demo endpoint (general limit: 200 req/min)
    print(f"TEST 1: /demo endpoint (limit: {DEMO_SUKSES_EXPECTED} req/min)")
    test_rate_limiting_sequential("/demo", num_requests=DEMO_TOTAL_REQUESTS, use_auth=False)
    
    print()
    
    # Test 2: Rate limiting on /api endpoint (API limit: 60 req/min)  
    print(f"TEST 2: /api/auth/me endpoint (limit: {API_SUKSES_EXPECTED} req/min)")
    test_rate_limiting_sequential("/api/auth/me", num_requests=API_TOTAL_REQUESTS, use_auth=True)
    
    print("\n🎉 Semua test selesai!")

if __name__ == "__main__":
    main()
