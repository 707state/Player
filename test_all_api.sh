#!/bin/bash

BASE_URL="http://192.168.246.1:8011"
TOTAL_TESTS=0
FAILED_TESTS=0

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试结果检查函数
check_response() {
    local expected_status=$1
    local response=$2
    local test_name=$3
    local http_code=$4
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ $test_name - Status code $http_code${NC}"
        # 检查响应内容
        if [ ! -z "$response" ]; then
            if echo "$response" | jq . >/dev/null 2>&1; then
                echo -e "${GREEN}✓ Response is valid JSON${NC}"
            else
                echo -e "${RED}✗ Invalid JSON response${NC}"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
        fi
    else
        echo -e "${RED}✗ $test_name - Expected status $expected_status but got $http_code${NC}"
        echo -e "${RED}Response: $response${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

test_endpoint() {
    local endpoint=$1
    local name=$2
    echo -e "\nTesting ${name} API..."
    echo "------------------------"

    # Test GET all
    echo "1. GET all ${name}:"
    response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}")
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 200 "$content" "GET all ${name}" "$http_code"

    # Test POST (Create)
    echo -e "\n2. POST new ${name}:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Album",
                    "artist": "Test Artist",
                    "genre": "Rock",
                    "year": 2023,
                    "cuts": ["Song 1", "Song 2"],
                    "url": "http://example.com",
                    "artwork": "http://example.com/cover.jpg",
                    "comment": "Test comment",
                    "rating": 5
                }')
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Book",
                    "author": "Test Author",
                    "genre": "Fiction",
                    "year": 2023,
                    "url": "http://example.com",
                    "cover": "http://example.com/cover.jpg",
                    "comment": "Test comment",
                    "rating": 5
                }')
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Movie",
                    "director": "Test Director",
                    "genre": "Action",
                    "year": 2023,
                    "url": "http://example.com",
                    "comment": "Test comment",
                    "rating": 5
                }')
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 201 "$content" "POST new ${name}" "$http_code"

    # Verify created item exists
    echo -e "\n3. Verify created item:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Album&artist=Test+Artist")
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Book&author=Test+Author")
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Movie&director=Test+Director")
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    if [ "$http_code" -eq 200 ] && echo "$content" | jq -e '. | length == 1' >/dev/null; then
        echo -e "${GREEN}✓ Item found after creation${NC}"
    else
        echo -e "${RED}✗ Item not found after creation${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    # Test UPDATE
    echo -e "\n4. UPDATE ${name}:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Album",
                    "artist": "Test Artist",
                    "genre": "Jazz",
                    "rating": 4
                }')
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Book",
                    "author": "Test Author",
                    "genre": "Non-Fiction",
                    "rating": 4
                }')
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Movie",
                    "director": "Test Director",
                    "genre": "Drama",
                    "rating": 4
                }')
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 200 "$content" "UPDATE ${name}" "$http_code"

    # Verify update
    echo -e "\n5. Verify update:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Album")
            http_code=$(echo "$response" | tail -n1)
            content=$(echo "$response" | sed '$ d')
            if [ "$http_code" -eq 200 ] && echo "$content" | jq -e '.[0].genre == "Jazz" and .[0].rating == 4' >/dev/null; then
                echo -e "${GREEN}✓ Update verified${NC}"
            else
                echo -e "${RED}✗ Update verification failed${NC}"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Book")
            http_code=$(echo "$response" | tail -n1)
            content=$(echo "$response" | sed '$ d')
            if [ "$http_code" -eq 200 ] && echo "$content" | jq -e '.[0].genre == "Non-Fiction" and .[0].rating == 4' >/dev/null; then
                echo -e "${GREEN}✓ Update verified${NC}"
            else
                echo -e "${RED}✗ Update verification failed${NC}"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Movie")
            http_code=$(echo "$response" | tail -n1)
            content=$(echo "$response" | sed '$ d')
            if [ "$http_code" -eq 200 ] && echo "$content" | jq -e '.[0].genre == "Drama" and .[0].rating == 4' >/dev/null; then
                echo -e "${GREEN}✓ Update verified${NC}"
            else
                echo -e "${RED}✗ Update verification failed${NC}"
                FAILED_TESTS=$((FAILED_TESTS + 1))
            fi
            ;;
    esac

    # Test DELETE
    echo -e "\n6. DELETE ${name}:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Album",
                    "artist": "Test Artist"
                }')
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Book",
                    "author": "Test Author"
                }')
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Test Movie",
                    "director": "Test Director"
                }')
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 200 "$content" "DELETE ${name}" "$http_code"

    # Verify deletion
    echo -e "\n7. Verify deletion:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Album")
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Book")
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X GET "${BASE_URL}${endpoint}?title=Test+Movie")
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    if [ "$http_code" -eq 200 ] && echo "$content" | jq -e '. | length == 0' >/dev/null; then
        echo -e "${GREEN}✓ Item successfully deleted${NC}"
    else
        echo -e "${RED}✗ Item still exists after deletion${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    # Test Error Cases
    echo -e "\n8. Error cases:"
    
    # Test POST with missing required fields
    echo "8.1. POST with missing required fields:"
    response=$(curl -s -w "\n%{http_code}" -X POST "${BASE_URL}${endpoint}" \
        -H "Content-Type: application/json" \
        -d '{"genre": "Test"}')
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 400 "$content" "POST with missing fields" "$http_code"

    # Test DELETE non-existent
    echo "8.2. DELETE non-existent item:"
    case "${endpoint}" in
        "/music")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Non Existent",
                    "artist": "Non Existent"
                }')
            ;;
        "/books")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Non Existent",
                    "author": "Non Existent"
                }')
            ;;
        "/movies")
            response=$(curl -s -w "\n%{http_code}" -X DELETE "${BASE_URL}${endpoint}" \
                -H "Content-Type: application/json" \
                -d '{
                    "title": "Non Existent",
                    "director": "Non Existent"
                }')
            ;;
    esac
    http_code=$(echo "$response" | tail -n1)
    content=$(echo "$response" | sed '$ d')
    check_response 404 "$content" "DELETE non-existent" "$http_code"
}

# 运行所有测试
test_endpoint "/music" "Music"
test_endpoint "/books" "Books"
test_endpoint "/movies" "Movies"

# 输出测试总结
echo -e "\n=== Test Summary ==="
echo "Total tests: $TOTAL_TESTS"
echo "Failed tests: $FAILED_TESTS"
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
