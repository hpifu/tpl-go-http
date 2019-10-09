Feature: echo 测试

    Scenario: echo
        When http 请求 GET /echo
            """
            {
                "params": {
                    "message": "hello world"
                }
            }
            """
        Then http 检查 200
            """
            {
                "json": {
                    "message": "hello world"
                }
            }
            """