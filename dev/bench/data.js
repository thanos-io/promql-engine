window.BENCHMARK_DATA = {
  "lastUpdate": 1675152463105,
  "repoUrl": "https://github.com/thanos-community/promql-engine",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "benye@amazon.com",
            "name": "Ben Ye",
            "username": "yeya24"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "075be8b8efd7d992700ae770c7126028bb82963a",
          "message": "Add continuous benchmark action for the new engine (#117)\n\n* add continuous benchmark action\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>\r\n\r\n* remove pr trigger\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>",
          "timestamp": "2022-11-10T08:42:00+01:00",
          "tree_id": "c2f14d11d1e2ffa2d19d3433442955ff7550301f",
          "url": "https://github.com/thanos-community/promql-engine/commit/075be8b8efd7d992700ae770c7126028bb82963a"
        },
        "date": 1668066347326,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 90607164,
            "unit": "ns/op\t28507064 B/op\t  126603 allocs/op",
            "extra": "15 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 83685365,
            "unit": "ns/op\t28690147 B/op\t  126620 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 83288999,
            "unit": "ns/op\t28652094 B/op\t  126616 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 83263237,
            "unit": "ns/op\t28637130 B/op\t  126612 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 82696248,
            "unit": "ns/op\t28641080 B/op\t  126614 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 74052484,
            "unit": "ns/op\t 9357205 B/op\t  121230 allocs/op",
            "extra": "16 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 72996408,
            "unit": "ns/op\t 9365085 B/op\t  121237 allocs/op",
            "extra": "15 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 75735980,
            "unit": "ns/op\t 9367476 B/op\t  121239 allocs/op",
            "extra": "16 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 76536613,
            "unit": "ns/op\t 9358708 B/op\t  121233 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 78393584,
            "unit": "ns/op\t 9427754 B/op\t  121241 allocs/op",
            "extra": "15 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 87759106,
            "unit": "ns/op\t18838061 B/op\t  206336 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 84488013,
            "unit": "ns/op\t18524776 B/op\t  206309 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 80909926,
            "unit": "ns/op\t18553466 B/op\t  206312 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 78372855,
            "unit": "ns/op\t18636582 B/op\t  206314 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 80461787,
            "unit": "ns/op\t18530427 B/op\t  206302 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 172141229,
            "unit": "ns/op\t30010712 B/op\t  150597 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 168815712,
            "unit": "ns/op\t30169733 B/op\t  150607 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 165911914,
            "unit": "ns/op\t30431227 B/op\t  150633 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 170037202,
            "unit": "ns/op\t30163411 B/op\t  150611 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 166416158,
            "unit": "ns/op\t30200006 B/op\t  150606 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 168324519,
            "unit": "ns/op\t11249897 B/op\t  145256 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 157299347,
            "unit": "ns/op\t11320276 B/op\t  145269 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 154200563,
            "unit": "ns/op\t11239328 B/op\t  145250 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 156847008,
            "unit": "ns/op\t11395478 B/op\t  145275 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 162435380,
            "unit": "ns/op\t11261581 B/op\t  145264 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 174250048,
            "unit": "ns/op\t20253628 B/op\t  230323 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 170864042,
            "unit": "ns/op\t20212916 B/op\t  230310 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 164516152,
            "unit": "ns/op\t20251912 B/op\t  230326 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 176143553,
            "unit": "ns/op\t20253594 B/op\t  230323 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 177152090,
            "unit": "ns/op\t20238636 B/op\t  230321 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 44550375,
            "unit": "ns/op\t14778209 B/op\t   98308 allocs/op",
            "extra": "25 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 43066960,
            "unit": "ns/op\t14793324 B/op\t   98318 allocs/op",
            "extra": "28 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 44659512,
            "unit": "ns/op\t14768235 B/op\t   98306 allocs/op",
            "extra": "25 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 43917566,
            "unit": "ns/op\t14770283 B/op\t   98308 allocs/op",
            "extra": "26 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 44918726,
            "unit": "ns/op\t14782277 B/op\t   98312 allocs/op",
            "extra": "26 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 105935433,
            "unit": "ns/op\t35072074 B/op\t  191932 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 106952454,
            "unit": "ns/op\t35036952 B/op\t  191921 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 108542326,
            "unit": "ns/op\t35141700 B/op\t  191954 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 108210446,
            "unit": "ns/op\t35101700 B/op\t  191969 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 103997397,
            "unit": "ns/op\t35039468 B/op\t  191935 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 89642889,
            "unit": "ns/op\t30861498 B/op\t  130577 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 89362995,
            "unit": "ns/op\t30805712 B/op\t  130567 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 87706889,
            "unit": "ns/op\t30843471 B/op\t  130570 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 88439591,
            "unit": "ns/op\t30827574 B/op\t  130562 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 86761672,
            "unit": "ns/op\t30828638 B/op\t  130569 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 88537034,
            "unit": "ns/op\t29989016 B/op\t  138950 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 89864243,
            "unit": "ns/op\t29990747 B/op\t  138947 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 89455699,
            "unit": "ns/op\t29990991 B/op\t  138950 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 87799724,
            "unit": "ns/op\t29975545 B/op\t  138945 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 88403524,
            "unit": "ns/op\t29996694 B/op\t  138950 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 93202495,
            "unit": "ns/op\t30488533 B/op\t  127561 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 93477796,
            "unit": "ns/op\t30542570 B/op\t  127576 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 87594949,
            "unit": "ns/op\t30482528 B/op\t  127556 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 86489706,
            "unit": "ns/op\t30499249 B/op\t  127557 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 89755646,
            "unit": "ns/op\t30479331 B/op\t  127549 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 75483291,
            "unit": "ns/op\t26887844 B/op\t   97820 allocs/op",
            "extra": "15 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 76047537,
            "unit": "ns/op\t26844998 B/op\t   97815 allocs/op",
            "extra": "16 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 78223212,
            "unit": "ns/op\t26852147 B/op\t   97817 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 79863062,
            "unit": "ns/op\t26843904 B/op\t   97813 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 77640946,
            "unit": "ns/op\t26834658 B/op\t   97809 allocs/op",
            "extra": "14 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 55018406,
            "unit": "ns/op\t35146689 B/op\t   75419 allocs/op",
            "extra": "22 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 50129574,
            "unit": "ns/op\t35147163 B/op\t   75420 allocs/op",
            "extra": "21 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 51875835,
            "unit": "ns/op\t35147714 B/op\t   75421 allocs/op",
            "extra": "24 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 51170024,
            "unit": "ns/op\t35146489 B/op\t   75418 allocs/op",
            "extra": "22 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 50245193,
            "unit": "ns/op\t35146819 B/op\t   75419 allocs/op",
            "extra": "21 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 48767472,
            "unit": "ns/op\t34957208 B/op\t   69421 allocs/op",
            "extra": "24 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 48625159,
            "unit": "ns/op\t34958350 B/op\t   69420 allocs/op",
            "extra": "25 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 46775422,
            "unit": "ns/op\t34957066 B/op\t   69420 allocs/op",
            "extra": "22 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 50218198,
            "unit": "ns/op\t34960285 B/op\t   69421 allocs/op",
            "extra": "25 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 48494150,
            "unit": "ns/op\t34958645 B/op\t   69421 allocs/op",
            "extra": "25 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 97101241,
            "unit": "ns/op\t29049876 B/op\t  130327 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 97436545,
            "unit": "ns/op\t29101162 B/op\t  130329 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 98189612,
            "unit": "ns/op\t29051492 B/op\t  130322 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 101274241,
            "unit": "ns/op\t29051734 B/op\t  130337 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 106144260,
            "unit": "ns/op\t29052615 B/op\t  130327 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 105334573,
            "unit": "ns/op\t29091496 B/op\t  129993 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 102009158,
            "unit": "ns/op\t29033092 B/op\t  129980 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 101385045,
            "unit": "ns/op\t29043477 B/op\t  129984 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 100770615,
            "unit": "ns/op\t29037333 B/op\t  129986 allocs/op",
            "extra": "13 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 100345127,
            "unit": "ns/op\t29104281 B/op\t  129995 allocs/op",
            "extra": "12 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 110533858,
            "unit": "ns/op\t31127189 B/op\t  135027 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 115552516,
            "unit": "ns/op\t31156940 B/op\t  135031 allocs/op",
            "extra": "9 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 111480132,
            "unit": "ns/op\t31129855 B/op\t  135021 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 111608992,
            "unit": "ns/op\t31131043 B/op\t  135023 allocs/op",
            "extra": "9 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 117587124,
            "unit": "ns/op\t31191423 B/op\t  135037 allocs/op",
            "extra": "10 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 158082239,
            "unit": "ns/op\t30156021 B/op\t  152053 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 160199610,
            "unit": "ns/op\t30238349 B/op\t  152061 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 159464674,
            "unit": "ns/op\t30239870 B/op\t  152065 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 159827206,
            "unit": "ns/op\t30247370 B/op\t  152070 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 152680503,
            "unit": "ns/op\t30312305 B/op\t  152068 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 172902985,
            "unit": "ns/op\t30217400 B/op\t  152050 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 177635927,
            "unit": "ns/op\t30172422 B/op\t  152057 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 175125741,
            "unit": "ns/op\t30272178 B/op\t  152077 allocs/op",
            "extra": "6 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 166003353,
            "unit": "ns/op\t30218893 B/op\t  152054 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 162300206,
            "unit": "ns/op\t30302765 B/op\t  152066 allocs/op",
            "extra": "7 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 301917590,
            "unit": "ns/op\t96749286 B/op\t  701247 allocs/op",
            "extra": "4 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 306212944,
            "unit": "ns/op\t97378054 B/op\t  701240 allocs/op",
            "extra": "4 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 311903444,
            "unit": "ns/op\t96971906 B/op\t  701230 allocs/op",
            "extra": "4 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 327200156,
            "unit": "ns/op\t96340118 B/op\t  701214 allocs/op",
            "extra": "4 times\n2 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 307777858,
            "unit": "ns/op\t96954848 B/op\t  701240 allocs/op",
            "extra": "4 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "437e914ef890465cb689b323e3540bb8d9ba3432",
          "message": ".github: use self-hosted runner (#118)\n\nI started a small runner on Equinix using CNCF's resources\r\n(c3.small.x86). Let's use it to have consistent results.",
          "timestamp": "2022-11-10T10:58:06+02:00",
          "tree_id": "e2ad364a66ee991f439f0dd266b1f98b51364171",
          "url": "https://github.com/thanos-community/promql-engine/commit/437e914ef890465cb689b323e3540bb8d9ba3432"
        },
        "date": 1668071121175,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23120225,
            "unit": "ns/op\t29488801 B/op\t  131575 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23457147,
            "unit": "ns/op\t29484337 B/op\t  131553 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23411400,
            "unit": "ns/op\t29498829 B/op\t  131572 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23510383,
            "unit": "ns/op\t29497642 B/op\t  131567 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24345525,
            "unit": "ns/op\t29504788 B/op\t  131566 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11015925,
            "unit": "ns/op\t12254050 B/op\t  126254 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11051543,
            "unit": "ns/op\t12215861 B/op\t  126231 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11032129,
            "unit": "ns/op\t12204306 B/op\t  126221 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11049268,
            "unit": "ns/op\t12225254 B/op\t  126221 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 10995732,
            "unit": "ns/op\t12201526 B/op\t  126233 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20434183,
            "unit": "ns/op\t21043885 B/op\t  211879 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20599917,
            "unit": "ns/op\t21064398 B/op\t  211888 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20481112,
            "unit": "ns/op\t21031938 B/op\t  211874 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20619206,
            "unit": "ns/op\t21043267 B/op\t  211880 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20639628,
            "unit": "ns/op\t21019515 B/op\t  211859 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29593339,
            "unit": "ns/op\t31304588 B/op\t  155333 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29659619,
            "unit": "ns/op\t31278144 B/op\t  155336 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29695246,
            "unit": "ns/op\t31320430 B/op\t  155379 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29528457,
            "unit": "ns/op\t31287905 B/op\t  155344 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29584446,
            "unit": "ns/op\t31291036 B/op\t  155327 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21657878,
            "unit": "ns/op\t13694550 B/op\t  149968 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21711615,
            "unit": "ns/op\t13775519 B/op\t  150055 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21725486,
            "unit": "ns/op\t13668291 B/op\t  149962 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21566924,
            "unit": "ns/op\t13687928 B/op\t  149969 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21683936,
            "unit": "ns/op\t13693237 B/op\t  150020 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30040968,
            "unit": "ns/op\t21907778 B/op\t  235327 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 29814989,
            "unit": "ns/op\t21937274 B/op\t  235357 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 29772861,
            "unit": "ns/op\t21857331 B/op\t  235292 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30086868,
            "unit": "ns/op\t21838980 B/op\t  235278 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 29986858,
            "unit": "ns/op\t21821218 B/op\t  235247 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18798030,
            "unit": "ns/op\t16013780 B/op\t  107820 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18833749,
            "unit": "ns/op\t16015921 B/op\t  107821 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18808770,
            "unit": "ns/op\t16022963 B/op\t  107840 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18718612,
            "unit": "ns/op\t16024458 B/op\t  107858 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18783393,
            "unit": "ns/op\t16039210 B/op\t  107885 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39876289,
            "unit": "ns/op\t36991008 B/op\t  201925 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39869145,
            "unit": "ns/op\t37024094 B/op\t  201951 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40194021,
            "unit": "ns/op\t37058511 B/op\t  201972 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39968489,
            "unit": "ns/op\t36973794 B/op\t  201908 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39875511,
            "unit": "ns/op\t37001806 B/op\t  201910 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33661697,
            "unit": "ns/op\t33110541 B/op\t  135933 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33652740,
            "unit": "ns/op\t33082099 B/op\t  135901 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33333358,
            "unit": "ns/op\t33113063 B/op\t  135925 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33933509,
            "unit": "ns/op\t33069324 B/op\t  135894 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33533124,
            "unit": "ns/op\t33066108 B/op\t  135900 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25863318,
            "unit": "ns/op\t30777038 B/op\t  143839 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25832540,
            "unit": "ns/op\t30791259 B/op\t  143853 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25691176,
            "unit": "ns/op\t30801793 B/op\t  143869 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25718918,
            "unit": "ns/op\t30774227 B/op\t  143848 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25749177,
            "unit": "ns/op\t30788227 B/op\t  143860 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33637429,
            "unit": "ns/op\t32762307 B/op\t  132908 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34047802,
            "unit": "ns/op\t32760843 B/op\t  132907 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33977428,
            "unit": "ns/op\t32673860 B/op\t  132864 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33776724,
            "unit": "ns/op\t32867111 B/op\t  132974 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33858184,
            "unit": "ns/op\t32786131 B/op\t  132914 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23901868,
            "unit": "ns/op\t27742519 B/op\t  102603 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23766700,
            "unit": "ns/op\t27755439 B/op\t  102617 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23642250,
            "unit": "ns/op\t27745230 B/op\t  102621 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23954648,
            "unit": "ns/op\t27755001 B/op\t  102630 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23734406,
            "unit": "ns/op\t27753439 B/op\t  102614 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21002522,
            "unit": "ns/op\t35236576 B/op\t   75784 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21053377,
            "unit": "ns/op\t35237572 B/op\t   75786 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21609766,
            "unit": "ns/op\t35237264 B/op\t   75786 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21545322,
            "unit": "ns/op\t35236661 B/op\t   75786 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21359420,
            "unit": "ns/op\t35236762 B/op\t   75785 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20755403,
            "unit": "ns/op\t35059434 B/op\t   69786 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20983246,
            "unit": "ns/op\t35061921 B/op\t   69787 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20897973,
            "unit": "ns/op\t35057679 B/op\t   69787 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20525502,
            "unit": "ns/op\t35054588 B/op\t   69785 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21055453,
            "unit": "ns/op\t35059278 B/op\t   69786 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41746183,
            "unit": "ns/op\t29871783 B/op\t  135385 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41912288,
            "unit": "ns/op\t29847896 B/op\t  135352 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42167788,
            "unit": "ns/op\t29864284 B/op\t  135402 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41680461,
            "unit": "ns/op\t29871951 B/op\t  135400 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42312623,
            "unit": "ns/op\t29876883 B/op\t  135369 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39354988,
            "unit": "ns/op\t29806071 B/op\t  134933 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38620371,
            "unit": "ns/op\t29794921 B/op\t  134928 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38583010,
            "unit": "ns/op\t29826015 B/op\t  134944 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38515376,
            "unit": "ns/op\t29825331 B/op\t  134948 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38819711,
            "unit": "ns/op\t29814630 B/op\t  134927 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50060373,
            "unit": "ns/op\t33351743 B/op\t  140395 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50539257,
            "unit": "ns/op\t33279884 B/op\t  140344 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50649788,
            "unit": "ns/op\t33344264 B/op\t  140378 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50401350,
            "unit": "ns/op\t33384753 B/op\t  140417 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50654603,
            "unit": "ns/op\t33304611 B/op\t  140374 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42314541,
            "unit": "ns/op\t31482978 B/op\t  157024 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42349805,
            "unit": "ns/op\t31488789 B/op\t  157046 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42417056,
            "unit": "ns/op\t31439482 B/op\t  156998 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42070619,
            "unit": "ns/op\t31488926 B/op\t  157049 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42431864,
            "unit": "ns/op\t31453563 B/op\t  157007 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42421780,
            "unit": "ns/op\t31488109 B/op\t  157042 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42244641,
            "unit": "ns/op\t31508149 B/op\t  157054 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42444444,
            "unit": "ns/op\t31481559 B/op\t  157048 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42433890,
            "unit": "ns/op\t31452712 B/op\t  156993 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42408842,
            "unit": "ns/op\t31441933 B/op\t  156995 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130255654,
            "unit": "ns/op\t98568157 B/op\t  704813 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 132319127,
            "unit": "ns/op\t98456460 B/op\t  704742 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131858084,
            "unit": "ns/op\t98567417 B/op\t  704814 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130696972,
            "unit": "ns/op\t98681143 B/op\t  704925 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131971556,
            "unit": "ns/op\t98495217 B/op\t  704751 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8afe40d8945907e7ba1f15b148c625803f2a22c3",
          "message": "Reduce allocations in aggregate (#119)\n\nThis commit reduces allocations in the hash aggregate operator by\r\nrecycling slices from the optional parameter operator and by reusing\r\nthe same slice to store parameter arguments.\r\n\r\nThe commit also adds a benchmark for this case.",
          "timestamp": "2022-11-10T14:30:41+01:00",
          "tree_id": "95c4a9ff7d7ed7143630b8fd689be9788182fe8f",
          "url": "https://github.com/thanos-community/promql-engine/commit/8afe40d8945907e7ba1f15b148c625803f2a22c3"
        },
        "date": 1668087192562,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23209066,
            "unit": "ns/op\t29516519 B/op\t  131402 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23219910,
            "unit": "ns/op\t29504644 B/op\t  131396 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23370618,
            "unit": "ns/op\t29481334 B/op\t  131368 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23372807,
            "unit": "ns/op\t29516441 B/op\t  131400 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24379719,
            "unit": "ns/op\t29471203 B/op\t  131357 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11038019,
            "unit": "ns/op\t12261302 B/op\t  126093 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11038012,
            "unit": "ns/op\t12268351 B/op\t  126065 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11066368,
            "unit": "ns/op\t12255452 B/op\t  126076 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11066573,
            "unit": "ns/op\t12273919 B/op\t  126082 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11018113,
            "unit": "ns/op\t12245209 B/op\t  126065 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20548224,
            "unit": "ns/op\t21038630 B/op\t  211704 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20529991,
            "unit": "ns/op\t21050698 B/op\t  211711 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20501861,
            "unit": "ns/op\t21038288 B/op\t  211708 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20387519,
            "unit": "ns/op\t21025350 B/op\t  211690 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20537734,
            "unit": "ns/op\t21034195 B/op\t  211700 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29546654,
            "unit": "ns/op\t31274451 B/op\t  155124 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29658469,
            "unit": "ns/op\t31296282 B/op\t  155156 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29547514,
            "unit": "ns/op\t31316336 B/op\t  155175 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29606108,
            "unit": "ns/op\t31300349 B/op\t  155158 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29710316,
            "unit": "ns/op\t31308281 B/op\t  155174 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21462382,
            "unit": "ns/op\t13780760 B/op\t  149860 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21646847,
            "unit": "ns/op\t13817309 B/op\t  149911 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21670458,
            "unit": "ns/op\t13694992 B/op\t  149846 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21846449,
            "unit": "ns/op\t13631534 B/op\t  149744 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21766774,
            "unit": "ns/op\t13667173 B/op\t  149752 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30118688,
            "unit": "ns/op\t21860489 B/op\t  235110 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30178872,
            "unit": "ns/op\t21891639 B/op\t  235127 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30130695,
            "unit": "ns/op\t21750241 B/op\t  235019 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30078008,
            "unit": "ns/op\t21827630 B/op\t  235083 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30113042,
            "unit": "ns/op\t21855540 B/op\t  235120 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29459171,
            "unit": "ns/op\t41194996 B/op\t  619334 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29298620,
            "unit": "ns/op\t41207706 B/op\t  619373 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29159550,
            "unit": "ns/op\t41206949 B/op\t  619322 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29216033,
            "unit": "ns/op\t41253877 B/op\t  619400 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29046794,
            "unit": "ns/op\t41211579 B/op\t  619363 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18863489,
            "unit": "ns/op\t16032408 B/op\t  107956 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18803457,
            "unit": "ns/op\t16036574 B/op\t  107973 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18874561,
            "unit": "ns/op\t16038405 B/op\t  107969 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18731193,
            "unit": "ns/op\t16029237 B/op\t  107954 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18726792,
            "unit": "ns/op\t16030982 B/op\t  107954 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39796332,
            "unit": "ns/op\t36990534 B/op\t  201794 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39817599,
            "unit": "ns/op\t36965738 B/op\t  201749 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39982181,
            "unit": "ns/op\t37007809 B/op\t  201794 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40102139,
            "unit": "ns/op\t37004731 B/op\t  201775 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39965940,
            "unit": "ns/op\t36977202 B/op\t  201763 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33386613,
            "unit": "ns/op\t33044804 B/op\t  135718 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33484841,
            "unit": "ns/op\t33055944 B/op\t  135726 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33709449,
            "unit": "ns/op\t33127007 B/op\t  135744 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33900583,
            "unit": "ns/op\t33067286 B/op\t  135715 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33614628,
            "unit": "ns/op\t33147146 B/op\t  135773 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25687997,
            "unit": "ns/op\t30793652 B/op\t  143659 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25705546,
            "unit": "ns/op\t30822777 B/op\t  143701 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25778546,
            "unit": "ns/op\t30790978 B/op\t  143666 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25727832,
            "unit": "ns/op\t30795236 B/op\t  143677 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25607461,
            "unit": "ns/op\t30780431 B/op\t  143680 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33942479,
            "unit": "ns/op\t32792620 B/op\t  132740 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34042585,
            "unit": "ns/op\t32802685 B/op\t  132746 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33833763,
            "unit": "ns/op\t32728980 B/op\t  132704 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33434656,
            "unit": "ns/op\t32813842 B/op\t  132754 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34043134,
            "unit": "ns/op\t32688564 B/op\t  132678 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23900372,
            "unit": "ns/op\t27722683 B/op\t  102471 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23832790,
            "unit": "ns/op\t27746667 B/op\t  102507 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23770118,
            "unit": "ns/op\t27737697 B/op\t  102497 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23821833,
            "unit": "ns/op\t27729303 B/op\t  102471 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 24096690,
            "unit": "ns/op\t27725010 B/op\t  102465 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21337487,
            "unit": "ns/op\t35235515 B/op\t   75724 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21482571,
            "unit": "ns/op\t35234878 B/op\t   75723 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21156268,
            "unit": "ns/op\t35235696 B/op\t   75724 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21613114,
            "unit": "ns/op\t35235858 B/op\t   75727 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21040213,
            "unit": "ns/op\t35235600 B/op\t   75725 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20764387,
            "unit": "ns/op\t35057608 B/op\t   69726 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20842545,
            "unit": "ns/op\t35058160 B/op\t   69726 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20692710,
            "unit": "ns/op\t35055942 B/op\t   69726 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20964421,
            "unit": "ns/op\t35061086 B/op\t   69726 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20560265,
            "unit": "ns/op\t35059048 B/op\t   69727 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42027555,
            "unit": "ns/op\t29865804 B/op\t  135213 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42180291,
            "unit": "ns/op\t29865432 B/op\t  135187 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41882914,
            "unit": "ns/op\t29880555 B/op\t  135214 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41980978,
            "unit": "ns/op\t29871517 B/op\t  135193 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41778589,
            "unit": "ns/op\t29838276 B/op\t  135166 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38664689,
            "unit": "ns/op\t29827266 B/op\t  134787 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38983229,
            "unit": "ns/op\t29823069 B/op\t  134759 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38878364,
            "unit": "ns/op\t29809017 B/op\t  134737 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38865296,
            "unit": "ns/op\t29792994 B/op\t  134729 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38873495,
            "unit": "ns/op\t29795179 B/op\t  134775 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50636070,
            "unit": "ns/op\t33396407 B/op\t  140240 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50832470,
            "unit": "ns/op\t33209053 B/op\t  140107 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51023541,
            "unit": "ns/op\t33246284 B/op\t  140138 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50331885,
            "unit": "ns/op\t33378393 B/op\t  140248 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50571067,
            "unit": "ns/op\t33239088 B/op\t  140128 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41981195,
            "unit": "ns/op\t31496533 B/op\t  156900 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42376702,
            "unit": "ns/op\t31462678 B/op\t  156856 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42139179,
            "unit": "ns/op\t31453974 B/op\t  156856 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42389771,
            "unit": "ns/op\t31491815 B/op\t  156863 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42532634,
            "unit": "ns/op\t31452115 B/op\t  156832 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42609362,
            "unit": "ns/op\t31472723 B/op\t  156855 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42481423,
            "unit": "ns/op\t31467606 B/op\t  156801 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42475973,
            "unit": "ns/op\t31464288 B/op\t  156835 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42337724,
            "unit": "ns/op\t31514739 B/op\t  156909 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42220618,
            "unit": "ns/op\t31456535 B/op\t  156849 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 132053746,
            "unit": "ns/op\t98638126 B/op\t  704833 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130308722,
            "unit": "ns/op\t98513548 B/op\t  704785 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131593135,
            "unit": "ns/op\t98516778 B/op\t  704702 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131040904,
            "unit": "ns/op\t98588241 B/op\t  704771 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130201042,
            "unit": "ns/op\t98562309 B/op\t  704775 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "approtas@amazon.com",
            "name": "Alan Protasio",
            "username": "alanprot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "39a919d5d4fcdc0052c19882ab4635891986f10b",
          "message": "Adding QueryType type (#121)",
          "timestamp": "2022-11-11T11:21:17-08:00",
          "tree_id": "de4e6108a6ce5ad05fcb316c501fc0529613cdbc",
          "url": "https://github.com/thanos-community/promql-engine/commit/39a919d5d4fcdc0052c19882ab4635891986f10b"
        },
        "date": 1668194628766,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23068517,
            "unit": "ns/op\t29516644 B/op\t  131479 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24194698,
            "unit": "ns/op\t29523233 B/op\t  131446 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24122675,
            "unit": "ns/op\t29545023 B/op\t  131458 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24245152,
            "unit": "ns/op\t29513545 B/op\t  131454 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 24295015,
            "unit": "ns/op\t29480508 B/op\t  131411 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11077954,
            "unit": "ns/op\t12266443 B/op\t  126121 allocs/op",
            "extra": "99 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 10992565,
            "unit": "ns/op\t12300104 B/op\t  126139 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11011291,
            "unit": "ns/op\t12260999 B/op\t  126118 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 10981057,
            "unit": "ns/op\t12243671 B/op\t  126096 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 10995598,
            "unit": "ns/op\t12294939 B/op\t  126127 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20588159,
            "unit": "ns/op\t21069339 B/op\t  211751 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20687450,
            "unit": "ns/op\t21070250 B/op\t  211749 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20649179,
            "unit": "ns/op\t21056321 B/op\t  211753 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20708080,
            "unit": "ns/op\t21080018 B/op\t  211763 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20593644,
            "unit": "ns/op\t21084299 B/op\t  211769 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29670532,
            "unit": "ns/op\t31333559 B/op\t  155231 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29764389,
            "unit": "ns/op\t31298904 B/op\t  155195 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29694269,
            "unit": "ns/op\t31300047 B/op\t  155219 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29830434,
            "unit": "ns/op\t31301479 B/op\t  155196 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29726760,
            "unit": "ns/op\t31316040 B/op\t  155216 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21658848,
            "unit": "ns/op\t13737751 B/op\t  149892 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21668504,
            "unit": "ns/op\t13746997 B/op\t  149889 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21688512,
            "unit": "ns/op\t13714584 B/op\t  149844 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21728961,
            "unit": "ns/op\t13697334 B/op\t  149829 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 21787460,
            "unit": "ns/op\t13633628 B/op\t  149803 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 29870019,
            "unit": "ns/op\t21873635 B/op\t  235153 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30046482,
            "unit": "ns/op\t21872138 B/op\t  235176 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30041556,
            "unit": "ns/op\t21871635 B/op\t  235171 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30098032,
            "unit": "ns/op\t21806848 B/op\t  235117 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 30051661,
            "unit": "ns/op\t21883258 B/op\t  235186 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29251959,
            "unit": "ns/op\t41225753 B/op\t  619488 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29188345,
            "unit": "ns/op\t41234125 B/op\t  619454 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29165448,
            "unit": "ns/op\t41286448 B/op\t  619472 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29285950,
            "unit": "ns/op\t41155408 B/op\t  619455 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29230578,
            "unit": "ns/op\t41279139 B/op\t  619471 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18793060,
            "unit": "ns/op\t16044951 B/op\t  108148 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18931443,
            "unit": "ns/op\t16049072 B/op\t  108161 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18877818,
            "unit": "ns/op\t16048182 B/op\t  108147 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18823793,
            "unit": "ns/op\t16049182 B/op\t  108162 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18821652,
            "unit": "ns/op\t16050231 B/op\t  108159 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40164900,
            "unit": "ns/op\t37027693 B/op\t  201905 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39769113,
            "unit": "ns/op\t37032822 B/op\t  201903 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40027242,
            "unit": "ns/op\t37013171 B/op\t  201891 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40096038,
            "unit": "ns/op\t36997375 B/op\t  201857 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39888796,
            "unit": "ns/op\t37024287 B/op\t  201874 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33606553,
            "unit": "ns/op\t33098279 B/op\t  135768 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33173175,
            "unit": "ns/op\t33195150 B/op\t  135843 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 34051417,
            "unit": "ns/op\t33110210 B/op\t  135788 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33632549,
            "unit": "ns/op\t33115688 B/op\t  135791 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33396326,
            "unit": "ns/op\t33128751 B/op\t  135799 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25801053,
            "unit": "ns/op\t30820452 B/op\t  143745 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25648045,
            "unit": "ns/op\t30828155 B/op\t  143748 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25857988,
            "unit": "ns/op\t30797004 B/op\t  143733 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25729446,
            "unit": "ns/op\t30820299 B/op\t  143742 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 25699316,
            "unit": "ns/op\t30800908 B/op\t  143736 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34074915,
            "unit": "ns/op\t32820689 B/op\t  132803 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33920025,
            "unit": "ns/op\t32691457 B/op\t  132738 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33634512,
            "unit": "ns/op\t32836465 B/op\t  132822 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33667197,
            "unit": "ns/op\t32699062 B/op\t  132740 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34275316,
            "unit": "ns/op\t32739218 B/op\t  132755 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23964129,
            "unit": "ns/op\t27748949 B/op\t  102522 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23878553,
            "unit": "ns/op\t27742405 B/op\t  102519 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23875124,
            "unit": "ns/op\t27764956 B/op\t  102526 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 24004567,
            "unit": "ns/op\t27695439 B/op\t  102481 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 24060709,
            "unit": "ns/op\t27736673 B/op\t  102526 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21589294,
            "unit": "ns/op\t35239318 B/op\t   75742 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21058183,
            "unit": "ns/op\t35239993 B/op\t   75743 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21088831,
            "unit": "ns/op\t35238908 B/op\t   75741 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21527768,
            "unit": "ns/op\t35239813 B/op\t   75743 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21547090,
            "unit": "ns/op\t35239996 B/op\t   75743 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20861773,
            "unit": "ns/op\t35061540 B/op\t   69744 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21100728,
            "unit": "ns/op\t35060341 B/op\t   69745 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20766839,
            "unit": "ns/op\t35059920 B/op\t   69744 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20888541,
            "unit": "ns/op\t35065356 B/op\t   69744 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20356389,
            "unit": "ns/op\t35061586 B/op\t   69743 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41766989,
            "unit": "ns/op\t29843717 B/op\t  135200 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41731593,
            "unit": "ns/op\t29844770 B/op\t  135217 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42184062,
            "unit": "ns/op\t29844396 B/op\t  135227 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42368133,
            "unit": "ns/op\t29879356 B/op\t  135251 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41779231,
            "unit": "ns/op\t29840677 B/op\t  135216 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39484108,
            "unit": "ns/op\t29836388 B/op\t  134833 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39061405,
            "unit": "ns/op\t29816258 B/op\t  134804 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39029242,
            "unit": "ns/op\t29812191 B/op\t  134791 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39165586,
            "unit": "ns/op\t29836670 B/op\t  134803 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38850836,
            "unit": "ns/op\t29817963 B/op\t  134794 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51007586,
            "unit": "ns/op\t33336605 B/op\t  140245 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50794365,
            "unit": "ns/op\t33297331 B/op\t  140231 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51201820,
            "unit": "ns/op\t33344408 B/op\t  140260 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50552360,
            "unit": "ns/op\t33219698 B/op\t  140175 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51026549,
            "unit": "ns/op\t33322440 B/op\t  140243 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42819315,
            "unit": "ns/op\t31392652 B/op\t  156788 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42065543,
            "unit": "ns/op\t31496020 B/op\t  156893 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42333735,
            "unit": "ns/op\t31443326 B/op\t  156863 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42237168,
            "unit": "ns/op\t31487418 B/op\t  156922 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42326283,
            "unit": "ns/op\t31468844 B/op\t  156909 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42258040,
            "unit": "ns/op\t31478961 B/op\t  156897 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42045640,
            "unit": "ns/op\t31492249 B/op\t  156919 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42304716,
            "unit": "ns/op\t31495117 B/op\t  156923 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42442875,
            "unit": "ns/op\t31484234 B/op\t  156894 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42257144,
            "unit": "ns/op\t31478124 B/op\t  156897 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 132527185,
            "unit": "ns/op\t98564791 B/op\t  704685 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131819903,
            "unit": "ns/op\t98597385 B/op\t  704708 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 132117843,
            "unit": "ns/op\t98614237 B/op\t  704764 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131612240,
            "unit": "ns/op\t98648445 B/op\t  704770 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129508447,
            "unit": "ns/op\t98631452 B/op\t  704767 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ac31e29297ca0b346753884bf791374ec4e7aa3d",
          "message": "Update prometheus/prometheus to v0.40.1 (#125)\n\nThis commit updates the prometheus/prometheus dependency to v0.40.1\r\nwhich adds support for native histograms in TSDB and PromQL.\r\nThe commit also unblocks updating this dependency in upstream Thanos.\r\n\r\nSince this engine does not yet support querying native histograms,\r\nthose cases will return an error for the time being. Unfortunately\r\nthe error has to be emitted during query execution since the\r\nhistogram_quantile function is now overloaded to work with both native\r\nand counter-based histograms.\r\n\r\nSupport for native-histograms in this engine will be added as a follow up.",
          "timestamp": "2022-11-12T23:34:07-08:00",
          "tree_id": "b2fa7d8326ec105d8924e6f815c67ccfff9ea006",
          "url": "https://github.com/thanos-community/promql-engine/commit/ac31e29297ca0b346753884bf791374ec4e7aa3d"
        },
        "date": 1668325007352,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25356722,
            "unit": "ns/op\t38773382 B/op\t  131445 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25364202,
            "unit": "ns/op\t38746728 B/op\t  131420 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25212847,
            "unit": "ns/op\t38773358 B/op\t  131433 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26384538,
            "unit": "ns/op\t38760857 B/op\t  131422 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26136165,
            "unit": "ns/op\t38774280 B/op\t  131430 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12026106,
            "unit": "ns/op\t12186382 B/op\t  126024 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12074407,
            "unit": "ns/op\t12190664 B/op\t  126046 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12060897,
            "unit": "ns/op\t12205163 B/op\t  126066 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12083606,
            "unit": "ns/op\t12203860 B/op\t  126042 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12089394,
            "unit": "ns/op\t12202842 B/op\t  126062 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21606034,
            "unit": "ns/op\t24264076 B/op\t  211832 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21460257,
            "unit": "ns/op\t24255485 B/op\t  211807 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21784010,
            "unit": "ns/op\t24242441 B/op\t  211817 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21330752,
            "unit": "ns/op\t24251025 B/op\t  211794 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21439032,
            "unit": "ns/op\t24274734 B/op\t  211829 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34976534,
            "unit": "ns/op\t41554862 B/op\t  155248 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34740626,
            "unit": "ns/op\t41530124 B/op\t  155217 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 35194551,
            "unit": "ns/op\t41594639 B/op\t  155273 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 35061166,
            "unit": "ns/op\t41581012 B/op\t  155274 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34539835,
            "unit": "ns/op\t41561946 B/op\t  155216 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24176451,
            "unit": "ns/op\t14680591 B/op\t  149853 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24185480,
            "unit": "ns/op\t14746726 B/op\t  149870 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24250874,
            "unit": "ns/op\t14636702 B/op\t  149785 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24057541,
            "unit": "ns/op\t14681945 B/op\t  149864 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24134585,
            "unit": "ns/op\t14778668 B/op\t  149897 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33663126,
            "unit": "ns/op\t26146004 B/op\t  235236 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33566949,
            "unit": "ns/op\t26103008 B/op\t  235193 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33706750,
            "unit": "ns/op\t26135233 B/op\t  235232 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33284655,
            "unit": "ns/op\t26109397 B/op\t  235204 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33593129,
            "unit": "ns/op\t26247801 B/op\t  235318 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30862427,
            "unit": "ns/op\t44484840 B/op\t  619491 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30778711,
            "unit": "ns/op\t44526965 B/op\t  619503 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30770198,
            "unit": "ns/op\t44495640 B/op\t  619498 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30778079,
            "unit": "ns/op\t44397410 B/op\t  619473 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30969005,
            "unit": "ns/op\t44538933 B/op\t  619520 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19838535,
            "unit": "ns/op\t19111008 B/op\t  107764 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20035586,
            "unit": "ns/op\t19112943 B/op\t  107758 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19914659,
            "unit": "ns/op\t19108303 B/op\t  107744 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19786223,
            "unit": "ns/op\t19107572 B/op\t  107746 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19873942,
            "unit": "ns/op\t19108309 B/op\t  107742 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43006130,
            "unit": "ns/op\t46266474 B/op\t  201781 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42277740,
            "unit": "ns/op\t46272839 B/op\t  201802 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42105036,
            "unit": "ns/op\t46254718 B/op\t  201792 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42614667,
            "unit": "ns/op\t46299055 B/op\t  201819 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42554828,
            "unit": "ns/op\t46257995 B/op\t  201793 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36674616,
            "unit": "ns/op\t42353800 B/op\t  135782 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36178090,
            "unit": "ns/op\t42264254 B/op\t  135746 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36030080,
            "unit": "ns/op\t42376618 B/op\t  135800 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36290597,
            "unit": "ns/op\t42292468 B/op\t  135754 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36067419,
            "unit": "ns/op\t42253669 B/op\t  135748 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28456605,
            "unit": "ns/op\t40052446 B/op\t  143722 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28122774,
            "unit": "ns/op\t40075368 B/op\t  143732 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28030961,
            "unit": "ns/op\t40075638 B/op\t  143734 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28074398,
            "unit": "ns/op\t40055472 B/op\t  143709 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28188478,
            "unit": "ns/op\t40103092 B/op\t  143755 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 35845216,
            "unit": "ns/op\t41963271 B/op\t  132756 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36090207,
            "unit": "ns/op\t41976985 B/op\t  132768 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36705809,
            "unit": "ns/op\t41869700 B/op\t  132713 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36421449,
            "unit": "ns/op\t41921220 B/op\t  132726 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36216461,
            "unit": "ns/op\t42095221 B/op\t  132820 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26100001,
            "unit": "ns/op\t37048996 B/op\t  102556 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 25915433,
            "unit": "ns/op\t37043086 B/op\t  102557 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26140040,
            "unit": "ns/op\t37073243 B/op\t  102567 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26510641,
            "unit": "ns/op\t37036717 B/op\t  102546 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26368107,
            "unit": "ns/op\t37005297 B/op\t  102520 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 24224568,
            "unit": "ns/op\t44501994 B/op\t   75742 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23108950,
            "unit": "ns/op\t44502621 B/op\t   75746 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 24028291,
            "unit": "ns/op\t44500917 B/op\t   75741 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23154421,
            "unit": "ns/op\t44501932 B/op\t   75743 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23229941,
            "unit": "ns/op\t44501584 B/op\t   75743 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22645493,
            "unit": "ns/op\t44324160 B/op\t   69746 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23009419,
            "unit": "ns/op\t44329902 B/op\t   69747 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21998443,
            "unit": "ns/op\t44321051 B/op\t   69745 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22322352,
            "unit": "ns/op\t44323650 B/op\t   69745 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21890856,
            "unit": "ns/op\t44325924 B/op\t   69745 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45152250,
            "unit": "ns/op\t39133813 B/op\t  135216 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45440847,
            "unit": "ns/op\t39135067 B/op\t  135241 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44905499,
            "unit": "ns/op\t39111885 B/op\t  135210 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44663274,
            "unit": "ns/op\t39102652 B/op\t  135175 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44055669,
            "unit": "ns/op\t39143180 B/op\t  135218 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42365361,
            "unit": "ns/op\t39097241 B/op\t  134816 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42485154,
            "unit": "ns/op\t39090667 B/op\t  134798 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42080701,
            "unit": "ns/op\t39094202 B/op\t  134779 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42079750,
            "unit": "ns/op\t39090354 B/op\t  134785 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42279638,
            "unit": "ns/op\t39095927 B/op\t  134791 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54950457,
            "unit": "ns/op\t42522668 B/op\t  140202 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54069987,
            "unit": "ns/op\t42575394 B/op\t  140241 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53299301,
            "unit": "ns/op\t42468437 B/op\t  140178 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53347677,
            "unit": "ns/op\t42601797 B/op\t  140271 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53440593,
            "unit": "ns/op\t42536429 B/op\t  140215 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46392432,
            "unit": "ns/op\t41630853 B/op\t  156914 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46399641,
            "unit": "ns/op\t41629343 B/op\t  156899 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46184275,
            "unit": "ns/op\t41638456 B/op\t  156904 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46072476,
            "unit": "ns/op\t41622405 B/op\t  156920 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46357975,
            "unit": "ns/op\t41628469 B/op\t  156907 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46027530,
            "unit": "ns/op\t41632979 B/op\t  156886 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46499929,
            "unit": "ns/op\t41626284 B/op\t  156890 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46238043,
            "unit": "ns/op\t41638222 B/op\t  156925 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46384390,
            "unit": "ns/op\t41610057 B/op\t  156900 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45678256,
            "unit": "ns/op\t41720423 B/op\t  157018 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131865566,
            "unit": "ns/op\t101847381 B/op\t  704948 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131951666,
            "unit": "ns/op\t101895930 B/op\t  705048 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 134389298,
            "unit": "ns/op\t101839538 B/op\t  704924 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 133173346,
            "unit": "ns/op\t101935586 B/op\t  705048 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131920488,
            "unit": "ns/op\t101861495 B/op\t  705004 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2af30177b51d15db8deae21adeaa2619083792fe",
          "message": "Fix histogram memory issue (#126)\n\nThis commit fixes the histogram_quantile memory issue by properly\r\nrecycling vector operator steps.",
          "timestamp": "2022-11-14T13:31:59+02:00",
          "tree_id": "a4a9311dbc18f520fc3037e0494619ed981c3e50",
          "url": "https://github.com/thanos-community/promql-engine/commit/2af30177b51d15db8deae21adeaa2619083792fe"
        },
        "date": 1668425677570,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25055990,
            "unit": "ns/op\t38783974 B/op\t  131159 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25155379,
            "unit": "ns/op\t38769548 B/op\t  131150 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25191776,
            "unit": "ns/op\t38739569 B/op\t  131100 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26131729,
            "unit": "ns/op\t38756755 B/op\t  131114 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26272636,
            "unit": "ns/op\t38796011 B/op\t  131138 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12154838,
            "unit": "ns/op\t12179363 B/op\t  125743 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12130299,
            "unit": "ns/op\t12201842 B/op\t  125749 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12087990,
            "unit": "ns/op\t12208346 B/op\t  125750 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12098357,
            "unit": "ns/op\t12134142 B/op\t  125694 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12181943,
            "unit": "ns/op\t12182246 B/op\t  125715 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21527908,
            "unit": "ns/op\t24221479 B/op\t  211487 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21373947,
            "unit": "ns/op\t24221797 B/op\t  211489 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21355244,
            "unit": "ns/op\t24217505 B/op\t  211490 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21845001,
            "unit": "ns/op\t24259227 B/op\t  211504 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21610883,
            "unit": "ns/op\t24229922 B/op\t  211481 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 35186706,
            "unit": "ns/op\t41499630 B/op\t  154880 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 35018490,
            "unit": "ns/op\t41533913 B/op\t  154924 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34897588,
            "unit": "ns/op\t41533326 B/op\t  154917 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34783527,
            "unit": "ns/op\t41550050 B/op\t  154944 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34959302,
            "unit": "ns/op\t41544848 B/op\t  154910 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23914854,
            "unit": "ns/op\t14728004 B/op\t  149629 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24236714,
            "unit": "ns/op\t14621046 B/op\t  149482 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24229306,
            "unit": "ns/op\t14650033 B/op\t  149511 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24152965,
            "unit": "ns/op\t14694731 B/op\t  149542 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24344751,
            "unit": "ns/op\t14590474 B/op\t  149458 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33485082,
            "unit": "ns/op\t26082844 B/op\t  234892 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33326634,
            "unit": "ns/op\t26028245 B/op\t  234838 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33458595,
            "unit": "ns/op\t25996245 B/op\t  234816 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33734435,
            "unit": "ns/op\t26151578 B/op\t  234924 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33545960,
            "unit": "ns/op\t26068299 B/op\t  234856 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30964823,
            "unit": "ns/op\t44467019 B/op\t  618863 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30852127,
            "unit": "ns/op\t44304574 B/op\t  618826 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30845766,
            "unit": "ns/op\t44429558 B/op\t  618834 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 31013222,
            "unit": "ns/op\t44444426 B/op\t  618856 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 31034053,
            "unit": "ns/op\t44337108 B/op\t  618806 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19827156,
            "unit": "ns/op\t19105655 B/op\t  107742 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19967750,
            "unit": "ns/op\t19121127 B/op\t  107762 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19959677,
            "unit": "ns/op\t19105765 B/op\t  107725 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19851074,
            "unit": "ns/op\t19111088 B/op\t  107735 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19880925,
            "unit": "ns/op\t19098670 B/op\t  107723 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42574176,
            "unit": "ns/op\t46218468 B/op\t  201437 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42594326,
            "unit": "ns/op\t46261868 B/op\t  201487 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43079283,
            "unit": "ns/op\t46229256 B/op\t  201508 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42945830,
            "unit": "ns/op\t46288810 B/op\t  201503 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42998110,
            "unit": "ns/op\t46270732 B/op\t  201514 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36588455,
            "unit": "ns/op\t42269353 B/op\t  135432 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 35630296,
            "unit": "ns/op\t42447176 B/op\t  135537 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 35796889,
            "unit": "ns/op\t42413384 B/op\t  135499 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36386227,
            "unit": "ns/op\t42223768 B/op\t  135433 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 35882312,
            "unit": "ns/op\t42355363 B/op\t  135481 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28202208,
            "unit": "ns/op\t40086552 B/op\t  143428 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28595912,
            "unit": "ns/op\t40003754 B/op\t  143371 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28657172,
            "unit": "ns/op\t40057398 B/op\t  143425 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28577707,
            "unit": "ns/op\t40027230 B/op\t  143374 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28242115,
            "unit": "ns/op\t40067585 B/op\t  143422 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36805542,
            "unit": "ns/op\t41933685 B/op\t  132434 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36452931,
            "unit": "ns/op\t41961423 B/op\t  132444 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 35918535,
            "unit": "ns/op\t41924615 B/op\t  132431 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 35910957,
            "unit": "ns/op\t41972821 B/op\t  132462 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 35649715,
            "unit": "ns/op\t41886944 B/op\t  132416 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26737178,
            "unit": "ns/op\t37014201 B/op\t  102325 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26382741,
            "unit": "ns/op\t37034295 B/op\t  102333 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26513736,
            "unit": "ns/op\t36966955 B/op\t  102290 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 25999657,
            "unit": "ns/op\t37052769 B/op\t  102345 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26478509,
            "unit": "ns/op\t36978686 B/op\t  102297 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23296774,
            "unit": "ns/op\t44492204 B/op\t   75632 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23809259,
            "unit": "ns/op\t44493044 B/op\t   75634 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23079391,
            "unit": "ns/op\t44493803 B/op\t   75634 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21758682,
            "unit": "ns/op\t44493220 B/op\t   75633 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23568889,
            "unit": "ns/op\t44493020 B/op\t   75632 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22264778,
            "unit": "ns/op\t44314870 B/op\t   69634 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22468814,
            "unit": "ns/op\t44315737 B/op\t   69636 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22712748,
            "unit": "ns/op\t44319714 B/op\t   69635 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23112084,
            "unit": "ns/op\t44319695 B/op\t   69635 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22530404,
            "unit": "ns/op\t44318147 B/op\t   69635 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44734925,
            "unit": "ns/op\t39119643 B/op\t  134936 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44927792,
            "unit": "ns/op\t39123304 B/op\t  134956 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45407132,
            "unit": "ns/op\t39083792 B/op\t  134843 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45778686,
            "unit": "ns/op\t39072038 B/op\t  134861 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44357414,
            "unit": "ns/op\t39122979 B/op\t  134947 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41428913,
            "unit": "ns/op\t39097982 B/op\t  134520 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42076336,
            "unit": "ns/op\t39044527 B/op\t  134461 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41610989,
            "unit": "ns/op\t39044356 B/op\t  134464 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42592992,
            "unit": "ns/op\t39047335 B/op\t  134484 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42189037,
            "unit": "ns/op\t39062390 B/op\t  134487 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54197069,
            "unit": "ns/op\t42453868 B/op\t  139846 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53453305,
            "unit": "ns/op\t42414108 B/op\t  139858 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52910272,
            "unit": "ns/op\t42579499 B/op\t  139928 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54294376,
            "unit": "ns/op\t42438548 B/op\t  139836 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53768943,
            "unit": "ns/op\t42492735 B/op\t  139868 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46344293,
            "unit": "ns/op\t41612323 B/op\t  156589 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46199736,
            "unit": "ns/op\t41599892 B/op\t  156552 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46237056,
            "unit": "ns/op\t41606749 B/op\t  156579 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46577457,
            "unit": "ns/op\t41579705 B/op\t  156552 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46162630,
            "unit": "ns/op\t41603856 B/op\t  156583 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46504848,
            "unit": "ns/op\t41585607 B/op\t  156573 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46397719,
            "unit": "ns/op\t41563443 B/op\t  156538 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46656774,
            "unit": "ns/op\t41581549 B/op\t  156580 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46081538,
            "unit": "ns/op\t41580493 B/op\t  156563 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46293706,
            "unit": "ns/op\t41576899 B/op\t  156553 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128162021,
            "unit": "ns/op\t60375514 B/op\t  704446 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130002083,
            "unit": "ns/op\t60603319 B/op\t  704535 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 127944179,
            "unit": "ns/op\t60444237 B/op\t  704465 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130284094,
            "unit": "ns/op\t60443504 B/op\t  704502 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129938426,
            "unit": "ns/op\t60484436 B/op\t  704502 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "approtas@amazon.com",
            "name": "Alan Protasio",
            "username": "alanprot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d68c8139c070320887b51d74c2bccd48a1d81eb0",
          "message": "Support topk and bottomK aggregations (#107)\n\n* Wip\r\n\r\n* Preseving the order returned by the topk operator\r\n\r\n* bench\r\n\r\n* opmization\r\n\r\n* lint\r\n\r\n* Fixing cases where the argument change values - range query\r\n\r\n* Adding test case where the arg vector selector yelds no result\r\n\r\n* comments\r\n\r\n* update readme\r\n\r\n* Inoring response order for TopK on instant query\r\n\r\n* Reduce allocations in aggregate\r\n\r\n* ADding the TODO on the test",
          "timestamp": "2022-11-15T11:24:26-08:00",
          "tree_id": "98ef72318778090dc70f9fcfe2a976cdd605af9e",
          "url": "https://github.com/thanos-community/promql-engine/commit/d68c8139c070320887b51d74c2bccd48a1d81eb0"
        },
        "date": 1668540443128,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25778716,
            "unit": "ns/op\t38789112 B/op\t  131431 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25990061,
            "unit": "ns/op\t38712985 B/op\t  131372 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26589654,
            "unit": "ns/op\t38722052 B/op\t  131356 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26888670,
            "unit": "ns/op\t38728210 B/op\t  131387 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26887340,
            "unit": "ns/op\t38725936 B/op\t  131378 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11917160,
            "unit": "ns/op\t12242598 B/op\t  126031 allocs/op",
            "extra": "99 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11878746,
            "unit": "ns/op\t12224256 B/op\t  126034 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11901231,
            "unit": "ns/op\t12197247 B/op\t  126027 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11877917,
            "unit": "ns/op\t12218940 B/op\t  126004 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11876208,
            "unit": "ns/op\t12210198 B/op\t  126028 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21244077,
            "unit": "ns/op\t24249632 B/op\t  211776 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21338694,
            "unit": "ns/op\t24282054 B/op\t  211784 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21188022,
            "unit": "ns/op\t24248177 B/op\t  211762 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21275598,
            "unit": "ns/op\t24260949 B/op\t  211789 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21580310,
            "unit": "ns/op\t24252345 B/op\t  211768 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16414135,
            "unit": "ns/op\t12131702 B/op\t  133594 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16344296,
            "unit": "ns/op\t12112647 B/op\t  133584 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16442618,
            "unit": "ns/op\t12102044 B/op\t  133533 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16408052,
            "unit": "ns/op\t12093965 B/op\t  133516 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16415055,
            "unit": "ns/op\t12134665 B/op\t  133594 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16410464,
            "unit": "ns/op\t12102134 B/op\t  132652 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16275886,
            "unit": "ns/op\t12048202 B/op\t  132570 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16299597,
            "unit": "ns/op\t12089446 B/op\t  132597 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16347892,
            "unit": "ns/op\t12043187 B/op\t  132576 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16367696,
            "unit": "ns/op\t12103765 B/op\t  132603 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33607159,
            "unit": "ns/op\t41553651 B/op\t  155220 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33671932,
            "unit": "ns/op\t41549541 B/op\t  155195 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33932279,
            "unit": "ns/op\t41563118 B/op\t  155216 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33504876,
            "unit": "ns/op\t41552752 B/op\t  155244 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33614136,
            "unit": "ns/op\t41557980 B/op\t  155226 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23906458,
            "unit": "ns/op\t14751517 B/op\t  149842 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23930625,
            "unit": "ns/op\t14693401 B/op\t  149831 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23961108,
            "unit": "ns/op\t14730817 B/op\t  149863 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24044543,
            "unit": "ns/op\t14653286 B/op\t  149803 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24219888,
            "unit": "ns/op\t14691487 B/op\t  149773 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33097193,
            "unit": "ns/op\t26093394 B/op\t  235150 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33113237,
            "unit": "ns/op\t26192892 B/op\t  235226 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32899501,
            "unit": "ns/op\t26122019 B/op\t  235165 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32831104,
            "unit": "ns/op\t26168476 B/op\t  235206 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32838261,
            "unit": "ns/op\t26197700 B/op\t  235238 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30406864,
            "unit": "ns/op\t44478309 B/op\t  619438 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30334957,
            "unit": "ns/op\t44412851 B/op\t  619422 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30352948,
            "unit": "ns/op\t44548042 B/op\t  619462 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30284111,
            "unit": "ns/op\t44466321 B/op\t  619450 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30110571,
            "unit": "ns/op\t44465202 B/op\t  619413 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20013095,
            "unit": "ns/op\t19114075 B/op\t  108178 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19973047,
            "unit": "ns/op\t19113336 B/op\t  108188 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19990501,
            "unit": "ns/op\t19126475 B/op\t  108216 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20082046,
            "unit": "ns/op\t19109776 B/op\t  108180 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19974023,
            "unit": "ns/op\t19123399 B/op\t  108212 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43265856,
            "unit": "ns/op\t46259734 B/op\t  202081 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42839775,
            "unit": "ns/op\t46275076 B/op\t  202092 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42457684,
            "unit": "ns/op\t46314013 B/op\t  202099 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42853387,
            "unit": "ns/op\t46286640 B/op\t  202119 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42748270,
            "unit": "ns/op\t46257934 B/op\t  202070 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36845921,
            "unit": "ns/op\t42400269 B/op\t  135777 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36115039,
            "unit": "ns/op\t42377265 B/op\t  135764 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36264047,
            "unit": "ns/op\t42250796 B/op\t  135698 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36463040,
            "unit": "ns/op\t42312652 B/op\t  135734 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36048879,
            "unit": "ns/op\t42334852 B/op\t  135751 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28943293,
            "unit": "ns/op\t40016412 B/op\t  143678 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28695196,
            "unit": "ns/op\t40010044 B/op\t  143668 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28603810,
            "unit": "ns/op\t40015509 B/op\t  143642 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28542807,
            "unit": "ns/op\t40039560 B/op\t  143690 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28877649,
            "unit": "ns/op\t40021542 B/op\t  143676 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36835190,
            "unit": "ns/op\t41964578 B/op\t  132734 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36576863,
            "unit": "ns/op\t41870088 B/op\t  132686 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36061966,
            "unit": "ns/op\t41879185 B/op\t  132697 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36635249,
            "unit": "ns/op\t41891805 B/op\t  132693 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36388628,
            "unit": "ns/op\t41898831 B/op\t  132696 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26662241,
            "unit": "ns/op\t37063307 B/op\t  102546 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26731841,
            "unit": "ns/op\t36996173 B/op\t  102476 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26877096,
            "unit": "ns/op\t36987203 B/op\t  102492 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 27187674,
            "unit": "ns/op\t36983000 B/op\t  102489 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26756234,
            "unit": "ns/op\t36966248 B/op\t  102476 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23267504,
            "unit": "ns/op\t44499220 B/op\t   75729 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23689874,
            "unit": "ns/op\t44500030 B/op\t   75730 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23281375,
            "unit": "ns/op\t44499026 B/op\t   75730 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23005548,
            "unit": "ns/op\t44499009 B/op\t   75730 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23737916,
            "unit": "ns/op\t44499061 B/op\t   75729 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22925770,
            "unit": "ns/op\t44321828 B/op\t   69731 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23664474,
            "unit": "ns/op\t44320245 B/op\t   69732 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21894943,
            "unit": "ns/op\t44323901 B/op\t   69732 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21904329,
            "unit": "ns/op\t44322288 B/op\t   69731 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21951539,
            "unit": "ns/op\t44316969 B/op\t   69730 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44706869,
            "unit": "ns/op\t39103648 B/op\t  135194 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44340010,
            "unit": "ns/op\t39145568 B/op\t  135218 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44640815,
            "unit": "ns/op\t39105620 B/op\t  135176 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44898433,
            "unit": "ns/op\t39095209 B/op\t  135172 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44273643,
            "unit": "ns/op\t39155728 B/op\t  135253 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42716470,
            "unit": "ns/op\t39053622 B/op\t  134722 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42039092,
            "unit": "ns/op\t39048672 B/op\t  134746 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41936481,
            "unit": "ns/op\t39089314 B/op\t  134754 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42460848,
            "unit": "ns/op\t39082564 B/op\t  134774 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42110576,
            "unit": "ns/op\t39069009 B/op\t  134770 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53161245,
            "unit": "ns/op\t42575168 B/op\t  140185 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53246769,
            "unit": "ns/op\t42501518 B/op\t  140159 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53828447,
            "unit": "ns/op\t42479530 B/op\t  140156 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53816551,
            "unit": "ns/op\t42582489 B/op\t  140214 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53285787,
            "unit": "ns/op\t42546499 B/op\t  140195 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46079894,
            "unit": "ns/op\t41614446 B/op\t  156860 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46105420,
            "unit": "ns/op\t41596819 B/op\t  156826 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45917084,
            "unit": "ns/op\t41565955 B/op\t  156809 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46343279,
            "unit": "ns/op\t41627075 B/op\t  156872 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45742500,
            "unit": "ns/op\t41641130 B/op\t  156902 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46289961,
            "unit": "ns/op\t41624407 B/op\t  156865 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45693964,
            "unit": "ns/op\t41615839 B/op\t  156872 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45551417,
            "unit": "ns/op\t41647456 B/op\t  156880 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45745066,
            "unit": "ns/op\t41643104 B/op\t  156912 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46496614,
            "unit": "ns/op\t41568655 B/op\t  156806 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129935220,
            "unit": "ns/op\t60717610 B/op\t  704456 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128556737,
            "unit": "ns/op\t60450071 B/op\t  704368 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129225810,
            "unit": "ns/op\t60675805 B/op\t  704469 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131223975,
            "unit": "ns/op\t60444569 B/op\t  704342 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129452992,
            "unit": "ns/op\t60393301 B/op\t  704304 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "44296a34973dd367f5946dbde2aacf8524bed00d",
          "message": "Fix join in binary-operator (#128)\n\n* Fix join in binary-operator\r\n\r\nThe binary operator assumes that only one low-cardinality series can\r\njoin any given high-cardinality series since PromQL does not allow many-to-many\r\njoins. However, this rule is true only within individual steps. Two low-cardinality\r\ntime-series can still join the same high-cardinality series if they do not have\r\noverlapping samples inside evaluation steps. This can happen during\r\nprocess restarts, where one time-series ends and new one starts.\r\n\r\nThis commit fixes the issue by properly calculating the output buffer for the\r\nlow-cardinality side in a binop. The case is hard to test right now because the\r\ntesting framework does not allow for defining offsets for time-series. I tested\r\nthis in one of our environments and I could not reproduce the original issue anymore.\r\n\r\n* Add test case",
          "timestamp": "2022-11-25T11:29:38+01:00",
          "tree_id": "88a284e88fefb5c60fcc1146b3cc78d7ddf79d51",
          "url": "https://github.com/thanos-community/promql-engine/commit/44296a34973dd367f5946dbde2aacf8524bed00d"
        },
        "date": 1669372343747,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25386715,
            "unit": "ns/op\t38759081 B/op\t  130896 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25660380,
            "unit": "ns/op\t38741729 B/op\t  130872 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26089369,
            "unit": "ns/op\t38706584 B/op\t  130823 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 27072436,
            "unit": "ns/op\t38731023 B/op\t  130854 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26626540,
            "unit": "ns/op\t38759437 B/op\t  130873 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12004113,
            "unit": "ns/op\t12237800 B/op\t  125514 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12002879,
            "unit": "ns/op\t12203714 B/op\t  125505 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12015783,
            "unit": "ns/op\t12212951 B/op\t  125492 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12001857,
            "unit": "ns/op\t12215487 B/op\t  125478 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12038976,
            "unit": "ns/op\t12197833 B/op\t  125495 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21172661,
            "unit": "ns/op\t24234797 B/op\t  211232 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21105595,
            "unit": "ns/op\t24250077 B/op\t  211250 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21111331,
            "unit": "ns/op\t24248657 B/op\t  211253 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21329655,
            "unit": "ns/op\t24214397 B/op\t  211236 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21008756,
            "unit": "ns/op\t24258041 B/op\t  211246 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16446790,
            "unit": "ns/op\t12029025 B/op\t  133908 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16589790,
            "unit": "ns/op\t12097468 B/op\t  133992 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16508519,
            "unit": "ns/op\t12073211 B/op\t  133913 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16444948,
            "unit": "ns/op\t11996385 B/op\t  133910 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16448155,
            "unit": "ns/op\t12032901 B/op\t  133895 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16485385,
            "unit": "ns/op\t12102086 B/op\t  134313 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16551826,
            "unit": "ns/op\t12089462 B/op\t  134234 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16508999,
            "unit": "ns/op\t12117589 B/op\t  134272 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16483068,
            "unit": "ns/op\t12057741 B/op\t  134229 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16430255,
            "unit": "ns/op\t12068680 B/op\t  134193 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33638902,
            "unit": "ns/op\t41544406 B/op\t  154678 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33662686,
            "unit": "ns/op\t41538548 B/op\t  154699 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33606535,
            "unit": "ns/op\t41516400 B/op\t  154658 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33557429,
            "unit": "ns/op\t41527295 B/op\t  154671 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33714473,
            "unit": "ns/op\t41546516 B/op\t  154701 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24397379,
            "unit": "ns/op\t14679624 B/op\t  149263 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24107860,
            "unit": "ns/op\t14724443 B/op\t  149309 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24374915,
            "unit": "ns/op\t14656100 B/op\t  149260 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24367470,
            "unit": "ns/op\t14661148 B/op\t  149256 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24388119,
            "unit": "ns/op\t14714520 B/op\t  149327 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32668840,
            "unit": "ns/op\t26143963 B/op\t  234662 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33340688,
            "unit": "ns/op\t26124878 B/op\t  234674 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33006504,
            "unit": "ns/op\t26069747 B/op\t  234636 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33288781,
            "unit": "ns/op\t26096966 B/op\t  234649 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32889098,
            "unit": "ns/op\t26049562 B/op\t  234617 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30220121,
            "unit": "ns/op\t44395367 B/op\t  618319 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30445616,
            "unit": "ns/op\t44341380 B/op\t  618300 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30279640,
            "unit": "ns/op\t44392664 B/op\t  618336 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30211869,
            "unit": "ns/op\t44418019 B/op\t  618351 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30234286,
            "unit": "ns/op\t44399523 B/op\t  618334 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19986081,
            "unit": "ns/op\t19121779 B/op\t  108084 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20070477,
            "unit": "ns/op\t19116244 B/op\t  108054 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20060449,
            "unit": "ns/op\t19122092 B/op\t  108078 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20070257,
            "unit": "ns/op\t19117145 B/op\t  108050 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19787312,
            "unit": "ns/op\t19128039 B/op\t  108070 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43321857,
            "unit": "ns/op\t46258195 B/op\t  201441 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43012596,
            "unit": "ns/op\t46250069 B/op\t  201433 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43099712,
            "unit": "ns/op\t46296046 B/op\t  201457 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42909975,
            "unit": "ns/op\t46288187 B/op\t  201460 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43157859,
            "unit": "ns/op\t46301768 B/op\t  201489 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36192856,
            "unit": "ns/op\t42307786 B/op\t  135189 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36346894,
            "unit": "ns/op\t42306213 B/op\t  135206 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37016623,
            "unit": "ns/op\t42261620 B/op\t  135171 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36501490,
            "unit": "ns/op\t42307520 B/op\t  135197 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36943449,
            "unit": "ns/op\t42212750 B/op\t  135148 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28660358,
            "unit": "ns/op\t39983172 B/op\t  143087 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28537338,
            "unit": "ns/op\t40017688 B/op\t  143143 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28536778,
            "unit": "ns/op\t39992230 B/op\t  143110 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28726607,
            "unit": "ns/op\t40008563 B/op\t  143125 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28803319,
            "unit": "ns/op\t40034874 B/op\t  143150 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36265530,
            "unit": "ns/op\t41916500 B/op\t  132179 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37231395,
            "unit": "ns/op\t41958331 B/op\t  132187 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36856218,
            "unit": "ns/op\t41934118 B/op\t  132199 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36085347,
            "unit": "ns/op\t41887342 B/op\t  132163 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36382191,
            "unit": "ns/op\t42012797 B/op\t  132226 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26798268,
            "unit": "ns/op\t36991205 B/op\t  102136 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26504082,
            "unit": "ns/op\t37016389 B/op\t  102148 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26775767,
            "unit": "ns/op\t36984148 B/op\t  102139 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26803346,
            "unit": "ns/op\t36997616 B/op\t  102137 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26817366,
            "unit": "ns/op\t36977096 B/op\t  102123 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23133688,
            "unit": "ns/op\t44492079 B/op\t   75550 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23777392,
            "unit": "ns/op\t44491619 B/op\t   75549 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23097592,
            "unit": "ns/op\t44491601 B/op\t   75548 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23542887,
            "unit": "ns/op\t44491051 B/op\t   75549 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23781640,
            "unit": "ns/op\t44491702 B/op\t   75549 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22264809,
            "unit": "ns/op\t44317974 B/op\t   69551 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22730663,
            "unit": "ns/op\t44311227 B/op\t   69550 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23052041,
            "unit": "ns/op\t44315463 B/op\t   69550 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22227313,
            "unit": "ns/op\t44319829 B/op\t   69553 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23350783,
            "unit": "ns/op\t44314873 B/op\t   69553 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44369999,
            "unit": "ns/op\t39100512 B/op\t  134658 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45005376,
            "unit": "ns/op\t39076404 B/op\t  134622 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44865881,
            "unit": "ns/op\t39070690 B/op\t  134628 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44273575,
            "unit": "ns/op\t39143063 B/op\t  134699 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45047408,
            "unit": "ns/op\t39092746 B/op\t  134650 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42838486,
            "unit": "ns/op\t39008551 B/op\t  134183 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42382003,
            "unit": "ns/op\t39058595 B/op\t  134216 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41479167,
            "unit": "ns/op\t39058462 B/op\t  134221 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41950049,
            "unit": "ns/op\t39060969 B/op\t  134220 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41796486,
            "unit": "ns/op\t39057679 B/op\t  134207 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53544686,
            "unit": "ns/op\t42393531 B/op\t  139569 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53599503,
            "unit": "ns/op\t42444418 B/op\t  139608 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53920164,
            "unit": "ns/op\t42480860 B/op\t  139601 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53531084,
            "unit": "ns/op\t42584459 B/op\t  139676 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53391908,
            "unit": "ns/op\t42557293 B/op\t  139668 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45693100,
            "unit": "ns/op\t41631696 B/op\t  156331 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45836409,
            "unit": "ns/op\t41588493 B/op\t  156337 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45848320,
            "unit": "ns/op\t41595583 B/op\t  156342 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45927281,
            "unit": "ns/op\t41601730 B/op\t  156330 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45896129,
            "unit": "ns/op\t41558582 B/op\t  156255 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46129066,
            "unit": "ns/op\t41561170 B/op\t  156275 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46633999,
            "unit": "ns/op\t41546916 B/op\t  156272 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45767963,
            "unit": "ns/op\t41616576 B/op\t  156345 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46096240,
            "unit": "ns/op\t41609989 B/op\t  156344 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45568963,
            "unit": "ns/op\t41592274 B/op\t  156314 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130423211,
            "unit": "ns/op\t60593790 B/op\t  704648 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129668322,
            "unit": "ns/op\t60580163 B/op\t  704683 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130430557,
            "unit": "ns/op\t60465205 B/op\t  704634 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130990733,
            "unit": "ns/op\t60723851 B/op\t  704688 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129147097,
            "unit": "ns/op\t60375733 B/op\t  704544 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "44ee950ce4b04963844e81c0e78aa4f24f0f578d",
          "message": "Implement matcher propagation (#115)\n\n* Implement matcher propagation\r\n\r\nThis commit implements matcher propagation between two vector selectors\r\nin a binary expression.\r\n\r\nThe feature is implemented as an opt-in optimizer which can be injected\r\nwhen instantiating the engine.\r\n\r\n* Use a single option to specify desired optimizers\r\n\r\n* Fix selecting logical optimizers\r\n\r\n* Rename Optimizers to AllOptimizers",
          "timestamp": "2022-11-28T18:11:36+01:00",
          "tree_id": "d88312af9ee796ea4016ac23918bd26c6660a701",
          "url": "https://github.com/thanos-community/promql-engine/commit/44ee950ce4b04963844e81c0e78aa4f24f0f578d"
        },
        "date": 1669655670527,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26168847,
            "unit": "ns/op\t38796655 B/op\t  131378 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26158347,
            "unit": "ns/op\t38797209 B/op\t  131355 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26800550,
            "unit": "ns/op\t38780954 B/op\t  131338 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26457666,
            "unit": "ns/op\t38793399 B/op\t  131352 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26853700,
            "unit": "ns/op\t38758142 B/op\t  131329 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12067814,
            "unit": "ns/op\t12196381 B/op\t  125934 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12040397,
            "unit": "ns/op\t12220392 B/op\t  125974 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12039074,
            "unit": "ns/op\t12204201 B/op\t  125938 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11973706,
            "unit": "ns/op\t12285637 B/op\t  126002 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12019154,
            "unit": "ns/op\t12236075 B/op\t  125973 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21298684,
            "unit": "ns/op\t24257492 B/op\t  211699 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21202701,
            "unit": "ns/op\t24246196 B/op\t  211706 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21155679,
            "unit": "ns/op\t24249769 B/op\t  211694 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21321894,
            "unit": "ns/op\t24258313 B/op\t  211717 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21216204,
            "unit": "ns/op\t24267157 B/op\t  211710 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16518934,
            "unit": "ns/op\t12149994 B/op\t  135578 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16501096,
            "unit": "ns/op\t12104092 B/op\t  135574 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16499292,
            "unit": "ns/op\t12110093 B/op\t  135567 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16488328,
            "unit": "ns/op\t12084834 B/op\t  135502 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16519704,
            "unit": "ns/op\t12123302 B/op\t  135571 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16544073,
            "unit": "ns/op\t12058268 B/op\t  131503 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16426908,
            "unit": "ns/op\t12010253 B/op\t  131541 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16475814,
            "unit": "ns/op\t12038556 B/op\t  131401 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16473984,
            "unit": "ns/op\t11997144 B/op\t  131425 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16448288,
            "unit": "ns/op\t12014924 B/op\t  131463 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33735270,
            "unit": "ns/op\t41566661 B/op\t  155137 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33992087,
            "unit": "ns/op\t41572490 B/op\t  155175 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33590401,
            "unit": "ns/op\t41578872 B/op\t  155185 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33923670,
            "unit": "ns/op\t41563336 B/op\t  155184 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33973162,
            "unit": "ns/op\t41585591 B/op\t  155167 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24361098,
            "unit": "ns/op\t14744132 B/op\t  149791 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24457362,
            "unit": "ns/op\t14678732 B/op\t  149736 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24387370,
            "unit": "ns/op\t14637238 B/op\t  149707 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24418836,
            "unit": "ns/op\t14643320 B/op\t  149699 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24251181,
            "unit": "ns/op\t14729148 B/op\t  149790 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33413328,
            "unit": "ns/op\t26121807 B/op\t  235101 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32919152,
            "unit": "ns/op\t26106028 B/op\t  235102 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33122230,
            "unit": "ns/op\t26097726 B/op\t  235111 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33102093,
            "unit": "ns/op\t26139124 B/op\t  235126 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33301904,
            "unit": "ns/op\t26125582 B/op\t  235120 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30354744,
            "unit": "ns/op\t44437432 B/op\t  619274 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30359562,
            "unit": "ns/op\t44492746 B/op\t  619289 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30513865,
            "unit": "ns/op\t44454178 B/op\t  619274 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30437911,
            "unit": "ns/op\t44441846 B/op\t  619255 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30122844,
            "unit": "ns/op\t44438888 B/op\t  619257 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19988841,
            "unit": "ns/op\t19124703 B/op\t  107946 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20064618,
            "unit": "ns/op\t19123329 B/op\t  107935 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20049774,
            "unit": "ns/op\t19129660 B/op\t  107956 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20042154,
            "unit": "ns/op\t19119601 B/op\t  107938 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20081948,
            "unit": "ns/op\t19114549 B/op\t  107928 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43030569,
            "unit": "ns/op\t46251372 B/op\t  201610 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42978725,
            "unit": "ns/op\t46301522 B/op\t  201691 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42692955,
            "unit": "ns/op\t46260960 B/op\t  201651 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42805727,
            "unit": "ns/op\t46288264 B/op\t  201688 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43085729,
            "unit": "ns/op\t46284292 B/op\t  201670 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36642113,
            "unit": "ns/op\t42331485 B/op\t  135681 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37204827,
            "unit": "ns/op\t42250710 B/op\t  135619 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36357524,
            "unit": "ns/op\t42392555 B/op\t  135717 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36725741,
            "unit": "ns/op\t42206368 B/op\t  135607 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36815590,
            "unit": "ns/op\t42272302 B/op\t  135653 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 29014285,
            "unit": "ns/op\t40019157 B/op\t  143586 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28785909,
            "unit": "ns/op\t40024109 B/op\t  143592 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28597781,
            "unit": "ns/op\t39985497 B/op\t  143560 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28504311,
            "unit": "ns/op\t40023165 B/op\t  143603 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28798951,
            "unit": "ns/op\t40024800 B/op\t  143588 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37040321,
            "unit": "ns/op\t41911310 B/op\t  132637 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37001130,
            "unit": "ns/op\t41975796 B/op\t  132660 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36918245,
            "unit": "ns/op\t41908153 B/op\t  132626 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37067777,
            "unit": "ns/op\t41920407 B/op\t  132654 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36687606,
            "unit": "ns/op\t41919243 B/op\t  132631 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26791552,
            "unit": "ns/op\t37030482 B/op\t  102467 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26748601,
            "unit": "ns/op\t36987205 B/op\t  102434 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26542717,
            "unit": "ns/op\t37016904 B/op\t  102475 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26939955,
            "unit": "ns/op\t36974159 B/op\t  102419 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26648615,
            "unit": "ns/op\t37006993 B/op\t  102440 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 24550205,
            "unit": "ns/op\t44498349 B/op\t   75707 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23305814,
            "unit": "ns/op\t44497555 B/op\t   75704 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23784364,
            "unit": "ns/op\t44498361 B/op\t   75707 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23637911,
            "unit": "ns/op\t44497909 B/op\t   75704 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23344981,
            "unit": "ns/op\t44497669 B/op\t   75705 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22838598,
            "unit": "ns/op\t44319756 B/op\t   69707 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22750379,
            "unit": "ns/op\t44321408 B/op\t   69706 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22368426,
            "unit": "ns/op\t44317148 B/op\t   69708 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22866001,
            "unit": "ns/op\t44323220 B/op\t   69708 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23953405,
            "unit": "ns/op\t44320203 B/op\t   69707 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44924990,
            "unit": "ns/op\t39099110 B/op\t  135079 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44959409,
            "unit": "ns/op\t39111533 B/op\t  135096 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44772573,
            "unit": "ns/op\t39104924 B/op\t  135110 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45091264,
            "unit": "ns/op\t39126992 B/op\t  135118 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44164309,
            "unit": "ns/op\t39093984 B/op\t  135119 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42298652,
            "unit": "ns/op\t39071832 B/op\t  134700 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42073561,
            "unit": "ns/op\t39066522 B/op\t  134670 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41538803,
            "unit": "ns/op\t39063641 B/op\t  134681 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41360129,
            "unit": "ns/op\t39080264 B/op\t  134692 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41928889,
            "unit": "ns/op\t39031449 B/op\t  134633 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54806694,
            "unit": "ns/op\t42520314 B/op\t  140103 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54092914,
            "unit": "ns/op\t42552951 B/op\t  140125 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53770889,
            "unit": "ns/op\t42480116 B/op\t  140098 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53725859,
            "unit": "ns/op\t42514979 B/op\t  140098 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54163107,
            "unit": "ns/op\t42536938 B/op\t  140122 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45393734,
            "unit": "ns/op\t41640903 B/op\t  156810 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45808686,
            "unit": "ns/op\t41623317 B/op\t  156790 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46018422,
            "unit": "ns/op\t41621772 B/op\t  156806 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45746712,
            "unit": "ns/op\t41615265 B/op\t  156762 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46223620,
            "unit": "ns/op\t41585581 B/op\t  156749 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46005672,
            "unit": "ns/op\t41612220 B/op\t  156771 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46156465,
            "unit": "ns/op\t41592925 B/op\t  156776 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45973712,
            "unit": "ns/op\t41584211 B/op\t  156749 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45387325,
            "unit": "ns/op\t41634315 B/op\t  156822 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45737124,
            "unit": "ns/op\t41606441 B/op\t  156755 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130104828,
            "unit": "ns/op\t60431377 B/op\t  704518 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128356952,
            "unit": "ns/op\t60566808 B/op\t  704609 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129715273,
            "unit": "ns/op\t60435455 B/op\t  704516 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130165737,
            "unit": "ns/op\t60512139 B/op\t  704574 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129944964,
            "unit": "ns/op\t60571080 B/op\t  704546 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f9757c8260ee4bde20d17308edc1ff3515ae0c41",
          "message": "Implement error reporting in binary operator for many-to-many joins (#131)\n\n* Implement error reporting in binary operator for many-to-many joins\r\n\r\nThe engine currently does not report many-to-many matching errors and\r\nwill likely silently return an incorrect result for such binary operations.\r\n\r\nThis commit adds support for this detection, similar to how the Prometheus\r\nengine reports these errors.\r\n\r\n* Remove error method from errManyToManyMatch",
          "timestamp": "2022-12-01T09:24:10+01:00",
          "tree_id": "cf8486a78263d0462a5a8d69bb020e5b73e72788",
          "url": "https://github.com/thanos-community/promql-engine/commit/f9757c8260ee4bde20d17308edc1ff3515ae0c41"
        },
        "date": 1669883219043,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25763974,
            "unit": "ns/op\t38796467 B/op\t  131637 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26107793,
            "unit": "ns/op\t38787345 B/op\t  131608 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26583801,
            "unit": "ns/op\t38820729 B/op\t  131630 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26616784,
            "unit": "ns/op\t38760340 B/op\t  131593 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26776371,
            "unit": "ns/op\t38791700 B/op\t  131627 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11878617,
            "unit": "ns/op\t12213488 B/op\t  126210 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11886307,
            "unit": "ns/op\t12228049 B/op\t  126211 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11802076,
            "unit": "ns/op\t12184277 B/op\t  126206 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11894242,
            "unit": "ns/op\t12207598 B/op\t  126222 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11849644,
            "unit": "ns/op\t12210286 B/op\t  126230 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21183034,
            "unit": "ns/op\t24255872 B/op\t  211968 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21257050,
            "unit": "ns/op\t24250734 B/op\t  211971 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21093489,
            "unit": "ns/op\t24279241 B/op\t  211975 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21127770,
            "unit": "ns/op\t24257604 B/op\t  211974 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21151341,
            "unit": "ns/op\t24257868 B/op\t  211977 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16371627,
            "unit": "ns/op\t12116931 B/op\t  133242 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16346290,
            "unit": "ns/op\t12101403 B/op\t  133245 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16336053,
            "unit": "ns/op\t12089919 B/op\t  133222 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16305013,
            "unit": "ns/op\t12106479 B/op\t  133221 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16314777,
            "unit": "ns/op\t12138984 B/op\t  133266 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16330787,
            "unit": "ns/op\t12066864 B/op\t  132993 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16394058,
            "unit": "ns/op\t12147022 B/op\t  133058 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16354709,
            "unit": "ns/op\t12128874 B/op\t  133081 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16405324,
            "unit": "ns/op\t12129583 B/op\t  133071 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16339472,
            "unit": "ns/op\t12137707 B/op\t  133111 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33551375,
            "unit": "ns/op\t41594075 B/op\t  155423 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33596077,
            "unit": "ns/op\t41578976 B/op\t  155440 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33437977,
            "unit": "ns/op\t41590297 B/op\t  155444 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33669305,
            "unit": "ns/op\t41638417 B/op\t  155470 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33652986,
            "unit": "ns/op\t41588264 B/op\t  155440 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24076170,
            "unit": "ns/op\t14772235 B/op\t  150062 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23890996,
            "unit": "ns/op\t14731226 B/op\t  150063 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24164594,
            "unit": "ns/op\t14702292 B/op\t  150008 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24060738,
            "unit": "ns/op\t14744717 B/op\t  150043 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24171332,
            "unit": "ns/op\t14729606 B/op\t  149999 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32886793,
            "unit": "ns/op\t26068536 B/op\t  235325 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33076194,
            "unit": "ns/op\t26138388 B/op\t  235381 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33068515,
            "unit": "ns/op\t26073296 B/op\t  235348 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33184911,
            "unit": "ns/op\t26147076 B/op\t  235375 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33099328,
            "unit": "ns/op\t26128194 B/op\t  235380 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30452439,
            "unit": "ns/op\t44494283 B/op\t  619885 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30194738,
            "unit": "ns/op\t44477972 B/op\t  619801 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30443739,
            "unit": "ns/op\t44459499 B/op\t  619846 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30301348,
            "unit": "ns/op\t44492125 B/op\t  619814 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30363154,
            "unit": "ns/op\t44521907 B/op\t  619864 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21087891,
            "unit": "ns/op\t22929600 B/op\t  108639 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21008157,
            "unit": "ns/op\t22934923 B/op\t  108653 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21006318,
            "unit": "ns/op\t22920802 B/op\t  108613 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21062145,
            "unit": "ns/op\t22938172 B/op\t  108667 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21081780,
            "unit": "ns/op\t22925806 B/op\t  108630 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46538771,
            "unit": "ns/op\t57813265 B/op\t  202798 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46570931,
            "unit": "ns/op\t57840378 B/op\t  202811 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46551000,
            "unit": "ns/op\t57860038 B/op\t  202863 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46838959,
            "unit": "ns/op\t57853871 B/op\t  202804 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46690193,
            "unit": "ns/op\t57855771 B/op\t  202843 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37249341,
            "unit": "ns/op\t42285751 B/op\t  135904 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36943920,
            "unit": "ns/op\t42344260 B/op\t  135940 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37189612,
            "unit": "ns/op\t42301895 B/op\t  135921 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36907210,
            "unit": "ns/op\t42332746 B/op\t  135942 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36372219,
            "unit": "ns/op\t42442921 B/op\t  135995 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28544935,
            "unit": "ns/op\t40031401 B/op\t  143870 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28662398,
            "unit": "ns/op\t40058374 B/op\t  143908 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28903871,
            "unit": "ns/op\t40049436 B/op\t  143901 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28659869,
            "unit": "ns/op\t40018873 B/op\t  143850 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28393928,
            "unit": "ns/op\t40038512 B/op\t  143887 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37028207,
            "unit": "ns/op\t41933436 B/op\t  132915 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36752393,
            "unit": "ns/op\t41970012 B/op\t  132928 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36850060,
            "unit": "ns/op\t41939865 B/op\t  132910 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37282745,
            "unit": "ns/op\t41875171 B/op\t  132881 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36878529,
            "unit": "ns/op\t42024181 B/op\t  132949 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26624178,
            "unit": "ns/op\t37002986 B/op\t  102621 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26539162,
            "unit": "ns/op\t36977442 B/op\t  102606 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26524331,
            "unit": "ns/op\t37020102 B/op\t  102625 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26461308,
            "unit": "ns/op\t37013012 B/op\t  102651 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26237349,
            "unit": "ns/op\t37023013 B/op\t  102648 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23252529,
            "unit": "ns/op\t44503494 B/op\t   75801 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23519434,
            "unit": "ns/op\t44504098 B/op\t   75801 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23315373,
            "unit": "ns/op\t44503804 B/op\t   75800 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23306457,
            "unit": "ns/op\t44504517 B/op\t   75802 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23127392,
            "unit": "ns/op\t44503860 B/op\t   75800 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22887620,
            "unit": "ns/op\t44329541 B/op\t   69803 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22873978,
            "unit": "ns/op\t44324634 B/op\t   69803 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22466434,
            "unit": "ns/op\t44329039 B/op\t   69804 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22386851,
            "unit": "ns/op\t44324339 B/op\t   69804 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22374050,
            "unit": "ns/op\t44328240 B/op\t   69804 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45290039,
            "unit": "ns/op\t39128880 B/op\t  135410 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44982499,
            "unit": "ns/op\t39140967 B/op\t  135403 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44665432,
            "unit": "ns/op\t39125792 B/op\t  135402 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44731545,
            "unit": "ns/op\t39132876 B/op\t  135388 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44854729,
            "unit": "ns/op\t39130970 B/op\t  135397 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42203975,
            "unit": "ns/op\t39087496 B/op\t  134952 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42153672,
            "unit": "ns/op\t39085996 B/op\t  134947 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41973593,
            "unit": "ns/op\t39093511 B/op\t  134978 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42426980,
            "unit": "ns/op\t39077826 B/op\t  134976 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42539189,
            "unit": "ns/op\t39054115 B/op\t  134935 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53291368,
            "unit": "ns/op\t42486387 B/op\t  140357 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54231979,
            "unit": "ns/op\t42670617 B/op\t  140450 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53191361,
            "unit": "ns/op\t42539238 B/op\t  140384 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53281904,
            "unit": "ns/op\t42529724 B/op\t  140388 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53807778,
            "unit": "ns/op\t42639657 B/op\t  140445 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46230357,
            "unit": "ns/op\t41638688 B/op\t  157060 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46022497,
            "unit": "ns/op\t41605726 B/op\t  157030 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45648982,
            "unit": "ns/op\t41622934 B/op\t  157044 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45995814,
            "unit": "ns/op\t41606928 B/op\t  157033 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45622408,
            "unit": "ns/op\t41615282 B/op\t  157015 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46095004,
            "unit": "ns/op\t41628805 B/op\t  157055 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45875703,
            "unit": "ns/op\t41588493 B/op\t  157006 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45414788,
            "unit": "ns/op\t41637912 B/op\t  157060 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45691046,
            "unit": "ns/op\t41625376 B/op\t  157049 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46016228,
            "unit": "ns/op\t41607263 B/op\t  157035 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128473515,
            "unit": "ns/op\t60314714 B/op\t  704235 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130598796,
            "unit": "ns/op\t60491142 B/op\t  704307 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130483492,
            "unit": "ns/op\t60431390 B/op\t  704293 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129921489,
            "unit": "ns/op\t60484924 B/op\t  704278 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128859735,
            "unit": "ns/op\t60401184 B/op\t  704225 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b4d311582a9e1d4a9bab4986ae8732011ccc58a6",
          "message": "execution: add support for KeepBool (#133)\n\nAdd support for \"KeepBool\"\r\nhttps://prometheus.io/docs/prometheus/latest/querying/operators/:\r\n\r\ntl;dr of the rules:\r\n- If `bool` modifier is provided between scalars then the result is either\r\n0 or 1\r\n- If `bool` modifier is provided between vectors/scalars then the result\r\nis 0 if the metric would've been dropped; `__name__` is removed.\r\n- If `bool` modifier is provided between vectors then the result\r\nis 0 if the metric would've been dropped; `__name__` is removed.",
          "timestamp": "2022-12-01T23:17:11-08:00",
          "tree_id": "0e6adde2fabe6a8fb36bfe8d690d42e5c039c8a3",
          "url": "https://github.com/thanos-community/promql-engine/commit/b4d311582a9e1d4a9bab4986ae8732011ccc58a6"
        },
        "date": 1669965601183,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25028999,
            "unit": "ns/op\t38725860 B/op\t  131079 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25630069,
            "unit": "ns/op\t38726906 B/op\t  131077 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25609722,
            "unit": "ns/op\t38729057 B/op\t  131052 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26369794,
            "unit": "ns/op\t38745924 B/op\t  131074 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26427311,
            "unit": "ns/op\t38747476 B/op\t  131076 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11876755,
            "unit": "ns/op\t12235478 B/op\t  125724 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11888741,
            "unit": "ns/op\t12188349 B/op\t  125697 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11916224,
            "unit": "ns/op\t12200105 B/op\t  125707 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11932802,
            "unit": "ns/op\t12219338 B/op\t  125729 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11899588,
            "unit": "ns/op\t12204506 B/op\t  125703 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21032869,
            "unit": "ns/op\t24229139 B/op\t  211462 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21201148,
            "unit": "ns/op\t24246921 B/op\t  211465 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21067063,
            "unit": "ns/op\t24248405 B/op\t  211459 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20908573,
            "unit": "ns/op\t24241064 B/op\t  211460 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20909566,
            "unit": "ns/op\t24237300 B/op\t  211460 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16303762,
            "unit": "ns/op\t12123266 B/op\t  134641 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16271978,
            "unit": "ns/op\t12087642 B/op\t  134569 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16291065,
            "unit": "ns/op\t12101221 B/op\t  134548 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16324143,
            "unit": "ns/op\t12111543 B/op\t  134638 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16351006,
            "unit": "ns/op\t12116997 B/op\t  134556 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16314068,
            "unit": "ns/op\t12061506 B/op\t  132977 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16335720,
            "unit": "ns/op\t12041491 B/op\t  132939 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16319828,
            "unit": "ns/op\t12096006 B/op\t  132998 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16355670,
            "unit": "ns/op\t12098763 B/op\t  132972 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16288675,
            "unit": "ns/op\t12076594 B/op\t  133002 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33420833,
            "unit": "ns/op\t41570669 B/op\t  154899 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33314881,
            "unit": "ns/op\t41558360 B/op\t  154923 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33123759,
            "unit": "ns/op\t41534622 B/op\t  154914 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33287840,
            "unit": "ns/op\t41549749 B/op\t  154929 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33367034,
            "unit": "ns/op\t41544704 B/op\t  154917 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23959405,
            "unit": "ns/op\t14674516 B/op\t  149492 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24060010,
            "unit": "ns/op\t14644990 B/op\t  149477 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23949598,
            "unit": "ns/op\t14712774 B/op\t  149496 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24051592,
            "unit": "ns/op\t14643228 B/op\t  149431 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23528380,
            "unit": "ns/op\t14716386 B/op\t  149548 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32907128,
            "unit": "ns/op\t26074331 B/op\t  234836 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32809194,
            "unit": "ns/op\t26052682 B/op\t  234811 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32864279,
            "unit": "ns/op\t26051323 B/op\t  234834 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32513813,
            "unit": "ns/op\t26050926 B/op\t  234815 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32990279,
            "unit": "ns/op\t26082663 B/op\t  234844 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30214645,
            "unit": "ns/op\t44452730 B/op\t  618811 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29982950,
            "unit": "ns/op\t44351081 B/op\t  618757 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30099542,
            "unit": "ns/op\t44448592 B/op\t  618803 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30304750,
            "unit": "ns/op\t44441350 B/op\t  618780 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 29928555,
            "unit": "ns/op\t44378631 B/op\t  618749 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21024866,
            "unit": "ns/op\t22925165 B/op\t  108617 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20983772,
            "unit": "ns/op\t22923324 B/op\t  108630 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21032838,
            "unit": "ns/op\t22933530 B/op\t  108626 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21004882,
            "unit": "ns/op\t22930311 B/op\t  108635 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 20990390,
            "unit": "ns/op\t22926547 B/op\t  108637 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46297576,
            "unit": "ns/op\t57794991 B/op\t  202358 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46430519,
            "unit": "ns/op\t57835576 B/op\t  202416 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46927724,
            "unit": "ns/op\t57867127 B/op\t  202463 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46225357,
            "unit": "ns/op\t57771891 B/op\t  202300 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 45935456,
            "unit": "ns/op\t57809851 B/op\t  202345 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37217508,
            "unit": "ns/op\t42291536 B/op\t  135405 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36652096,
            "unit": "ns/op\t42213065 B/op\t  135362 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36502842,
            "unit": "ns/op\t42339912 B/op\t  135433 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36768300,
            "unit": "ns/op\t42253212 B/op\t  135386 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36922316,
            "unit": "ns/op\t42230928 B/op\t  135373 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28429628,
            "unit": "ns/op\t40054018 B/op\t  143380 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28343543,
            "unit": "ns/op\t40032289 B/op\t  143370 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28518591,
            "unit": "ns/op\t40022698 B/op\t  143355 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28532146,
            "unit": "ns/op\t40034665 B/op\t  143375 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28743089,
            "unit": "ns/op\t39977840 B/op\t  143312 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36939643,
            "unit": "ns/op\t41853298 B/op\t  132355 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36469851,
            "unit": "ns/op\t41994751 B/op\t  132441 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37228445,
            "unit": "ns/op\t41836613 B/op\t  132347 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36836141,
            "unit": "ns/op\t41863053 B/op\t  132365 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36703887,
            "unit": "ns/op\t42010180 B/op\t  132443 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26458382,
            "unit": "ns/op\t37021440 B/op\t  102300 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26438855,
            "unit": "ns/op\t36972590 B/op\t  102270 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26643972,
            "unit": "ns/op\t36994027 B/op\t  102286 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26841942,
            "unit": "ns/op\t36966402 B/op\t  102269 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26418608,
            "unit": "ns/op\t37014139 B/op\t  102287 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23668696,
            "unit": "ns/op\t44492664 B/op\t   75621 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23273509,
            "unit": "ns/op\t44493485 B/op\t   75622 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 24042162,
            "unit": "ns/op\t44493275 B/op\t   75622 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23080423,
            "unit": "ns/op\t44493025 B/op\t   75622 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23461302,
            "unit": "ns/op\t44492783 B/op\t   75621 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22883053,
            "unit": "ns/op\t44317800 B/op\t   69625 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22715688,
            "unit": "ns/op\t44317346 B/op\t   69624 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22653709,
            "unit": "ns/op\t44316320 B/op\t   69624 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22005099,
            "unit": "ns/op\t44315261 B/op\t   69622 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22163926,
            "unit": "ns/op\t44316273 B/op\t   69623 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43864942,
            "unit": "ns/op\t39131076 B/op\t  134894 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44396923,
            "unit": "ns/op\t39097832 B/op\t  134842 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44705763,
            "unit": "ns/op\t39087443 B/op\t  134853 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44301529,
            "unit": "ns/op\t39105934 B/op\t  134867 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44873134,
            "unit": "ns/op\t39095887 B/op\t  134834 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41453245,
            "unit": "ns/op\t39099859 B/op\t  134480 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41385824,
            "unit": "ns/op\t39056195 B/op\t  134423 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41743710,
            "unit": "ns/op\t39056633 B/op\t  134438 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42145092,
            "unit": "ns/op\t39037160 B/op\t  134393 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41340573,
            "unit": "ns/op\t39080390 B/op\t  134459 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54129862,
            "unit": "ns/op\t42545580 B/op\t  139872 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54654239,
            "unit": "ns/op\t42430777 B/op\t  139822 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54532512,
            "unit": "ns/op\t42411529 B/op\t  139798 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53656967,
            "unit": "ns/op\t42615338 B/op\t  139923 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53875789,
            "unit": "ns/op\t42516754 B/op\t  139876 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46119056,
            "unit": "ns/op\t41593688 B/op\t  156539 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46051489,
            "unit": "ns/op\t41624268 B/op\t  156562 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45867695,
            "unit": "ns/op\t41615485 B/op\t  156547 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45928542,
            "unit": "ns/op\t41594875 B/op\t  156526 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45336452,
            "unit": "ns/op\t41613771 B/op\t  156536 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45884643,
            "unit": "ns/op\t41558578 B/op\t  156502 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45505262,
            "unit": "ns/op\t41603518 B/op\t  156558 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45478157,
            "unit": "ns/op\t41607459 B/op\t  156533 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45730790,
            "unit": "ns/op\t41571187 B/op\t  156505 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45761512,
            "unit": "ns/op\t41522178 B/op\t  156456 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130957646,
            "unit": "ns/op\t60483012 B/op\t  704450 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128349978,
            "unit": "ns/op\t60482263 B/op\t  704423 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128916339,
            "unit": "ns/op\t60511024 B/op\t  704460 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128345999,
            "unit": "ns/op\t60468255 B/op\t  704473 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128854292,
            "unit": "ns/op\t60363946 B/op\t  704420 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "saswataminsta@yahoo.com",
            "name": "Saswata Mukherjee",
            "username": "saswatamcode"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fa5e576beae29b2b4a2be5b652c8e9321afb4d57",
          "message": "Use efficientgo/core/errors consistently (#134)\n\nSigned-off-by: Saswata Mukherjee <saswataminsta@yahoo.com>\r\n\r\nSigned-off-by: Saswata Mukherjee <saswataminsta@yahoo.com>",
          "timestamp": "2022-12-02T12:05:31+01:00",
          "tree_id": "f23f9350b84f5efba7e72b6215ae402988df435e",
          "url": "https://github.com/thanos-community/promql-engine/commit/fa5e576beae29b2b4a2be5b652c8e9321afb4d57"
        },
        "date": 1669979300465,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25430174,
            "unit": "ns/op\t38814059 B/op\t  131233 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26348289,
            "unit": "ns/op\t38804913 B/op\t  131227 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26422356,
            "unit": "ns/op\t38793588 B/op\t  131193 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26333832,
            "unit": "ns/op\t38820885 B/op\t  131228 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26450197,
            "unit": "ns/op\t38834660 B/op\t  131227 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12274675,
            "unit": "ns/op\t12236930 B/op\t  125816 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12350612,
            "unit": "ns/op\t12203193 B/op\t  125816 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12372347,
            "unit": "ns/op\t12250361 B/op\t  125825 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12317338,
            "unit": "ns/op\t12182281 B/op\t  125791 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12310397,
            "unit": "ns/op\t12213916 B/op\t  125822 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21446538,
            "unit": "ns/op\t24255210 B/op\t  211569 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21418162,
            "unit": "ns/op\t24281319 B/op\t  211597 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21611695,
            "unit": "ns/op\t24253608 B/op\t  211575 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21526923,
            "unit": "ns/op\t24276719 B/op\t  211575 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21508389,
            "unit": "ns/op\t24292364 B/op\t  211588 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16724565,
            "unit": "ns/op\t12051121 B/op\t  133808 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16588033,
            "unit": "ns/op\t11987803 B/op\t  133697 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16616167,
            "unit": "ns/op\t12016791 B/op\t  133730 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16647509,
            "unit": "ns/op\t12029989 B/op\t  133745 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16612791,
            "unit": "ns/op\t12052969 B/op\t  133755 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16776772,
            "unit": "ns/op\t12042265 B/op\t  134184 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16638131,
            "unit": "ns/op\t12068309 B/op\t  134225 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16629078,
            "unit": "ns/op\t12046603 B/op\t  134162 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16609443,
            "unit": "ns/op\t12138941 B/op\t  134253 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16646354,
            "unit": "ns/op\t12047472 B/op\t  134184 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33886216,
            "unit": "ns/op\t41548567 B/op\t  154993 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34289112,
            "unit": "ns/op\t41543910 B/op\t  154992 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34033349,
            "unit": "ns/op\t41544270 B/op\t  155015 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34179158,
            "unit": "ns/op\t41563246 B/op\t  155015 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 34012451,
            "unit": "ns/op\t41571069 B/op\t  155020 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24407136,
            "unit": "ns/op\t14713472 B/op\t  149609 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24674959,
            "unit": "ns/op\t14578814 B/op\t  149506 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24377070,
            "unit": "ns/op\t14708241 B/op\t  149630 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24649972,
            "unit": "ns/op\t14708784 B/op\t  149618 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24309092,
            "unit": "ns/op\t14767064 B/op\t  149713 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33514852,
            "unit": "ns/op\t26102852 B/op\t  234945 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33417406,
            "unit": "ns/op\t26062108 B/op\t  234919 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33454687,
            "unit": "ns/op\t26087909 B/op\t  234921 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32911340,
            "unit": "ns/op\t26074063 B/op\t  234946 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33475171,
            "unit": "ns/op\t26148001 B/op\t  234992 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30637035,
            "unit": "ns/op\t44401084 B/op\t  618941 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30556802,
            "unit": "ns/op\t44420970 B/op\t  618972 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30484618,
            "unit": "ns/op\t44459741 B/op\t  618996 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30788578,
            "unit": "ns/op\t44413244 B/op\t  618977 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30564510,
            "unit": "ns/op\t44392134 B/op\t  618963 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21278183,
            "unit": "ns/op\t22924190 B/op\t  108498 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21306140,
            "unit": "ns/op\t22915412 B/op\t  108479 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21193271,
            "unit": "ns/op\t22947005 B/op\t  108533 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21317867,
            "unit": "ns/op\t22917525 B/op\t  108475 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21290938,
            "unit": "ns/op\t22926108 B/op\t  108497 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 47222487,
            "unit": "ns/op\t57862248 B/op\t  202435 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46793656,
            "unit": "ns/op\t57863511 B/op\t  202412 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46938498,
            "unit": "ns/op\t57833895 B/op\t  202378 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46907472,
            "unit": "ns/op\t57819950 B/op\t  202378 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 47036874,
            "unit": "ns/op\t57864060 B/op\t  202414 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37178242,
            "unit": "ns/op\t42319012 B/op\t  135525 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37169141,
            "unit": "ns/op\t42339077 B/op\t  135520 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36835519,
            "unit": "ns/op\t42259622 B/op\t  135496 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37188584,
            "unit": "ns/op\t42317246 B/op\t  135521 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36966296,
            "unit": "ns/op\t42377582 B/op\t  135542 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28904844,
            "unit": "ns/op\t40011400 B/op\t  143442 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28920889,
            "unit": "ns/op\t40015057 B/op\t  143434 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28924441,
            "unit": "ns/op\t40020083 B/op\t  143451 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28810918,
            "unit": "ns/op\t40026449 B/op\t  143452 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28936684,
            "unit": "ns/op\t39980065 B/op\t  143407 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37079509,
            "unit": "ns/op\t41912677 B/op\t  132495 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36899451,
            "unit": "ns/op\t41980706 B/op\t  132528 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37034316,
            "unit": "ns/op\t41924374 B/op\t  132501 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37097127,
            "unit": "ns/op\t41864671 B/op\t  132468 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37221258,
            "unit": "ns/op\t41890281 B/op\t  132480 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26734556,
            "unit": "ns/op\t37017142 B/op\t  102382 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26993759,
            "unit": "ns/op\t36975249 B/op\t  102342 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26869689,
            "unit": "ns/op\t36997880 B/op\t  102353 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 27059049,
            "unit": "ns/op\t36972926 B/op\t  102324 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26955822,
            "unit": "ns/op\t36983394 B/op\t  102347 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 24011785,
            "unit": "ns/op\t44497004 B/op\t   75658 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23045424,
            "unit": "ns/op\t44497195 B/op\t   75658 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23273182,
            "unit": "ns/op\t44497082 B/op\t   75658 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23027347,
            "unit": "ns/op\t44496716 B/op\t   75657 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 23844775,
            "unit": "ns/op\t44496794 B/op\t   75658 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22400264,
            "unit": "ns/op\t44315479 B/op\t   69659 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22486271,
            "unit": "ns/op\t44319492 B/op\t   69658 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 23220320,
            "unit": "ns/op\t44316499 B/op\t   69659 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22710269,
            "unit": "ns/op\t44321086 B/op\t   69659 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22631849,
            "unit": "ns/op\t44322710 B/op\t   69660 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44491647,
            "unit": "ns/op\t39141373 B/op\t  134980 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44493684,
            "unit": "ns/op\t39124596 B/op\t  134979 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44865788,
            "unit": "ns/op\t39111880 B/op\t  134961 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45054638,
            "unit": "ns/op\t39117224 B/op\t  134979 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44925286,
            "unit": "ns/op\t39128979 B/op\t  135000 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42019099,
            "unit": "ns/op\t39070316 B/op\t  134557 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41982503,
            "unit": "ns/op\t39089977 B/op\t  134549 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41583465,
            "unit": "ns/op\t39100769 B/op\t  134572 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41658131,
            "unit": "ns/op\t39068032 B/op\t  134522 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42406664,
            "unit": "ns/op\t39050904 B/op\t  134529 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54808594,
            "unit": "ns/op\t42465209 B/op\t  139936 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54070560,
            "unit": "ns/op\t42633187 B/op\t  140041 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54275054,
            "unit": "ns/op\t42590923 B/op\t  140016 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54545539,
            "unit": "ns/op\t42536014 B/op\t  139972 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54436666,
            "unit": "ns/op\t42554981 B/op\t  139978 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46075001,
            "unit": "ns/op\t41611822 B/op\t  156637 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46116086,
            "unit": "ns/op\t41628960 B/op\t  156651 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46229597,
            "unit": "ns/op\t41644521 B/op\t  156661 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46405389,
            "unit": "ns/op\t41595763 B/op\t  156590 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 45694397,
            "unit": "ns/op\t41601268 B/op\t  156655 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46172817,
            "unit": "ns/op\t41610997 B/op\t  156668 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45974410,
            "unit": "ns/op\t41640343 B/op\t  156655 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46124781,
            "unit": "ns/op\t41609946 B/op\t  156632 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 45823922,
            "unit": "ns/op\t41631992 B/op\t  156671 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46357387,
            "unit": "ns/op\t41577890 B/op\t  156590 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128737334,
            "unit": "ns/op\t60568838 B/op\t  704571 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128550486,
            "unit": "ns/op\t60297799 B/op\t  704437 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130029040,
            "unit": "ns/op\t60433986 B/op\t  704510 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129264751,
            "unit": "ns/op\t60416096 B/op\t  704463 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130046868,
            "unit": "ns/op\t60411614 B/op\t  704512 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fa1f96080ee7c732ad8f290db83ad7559aa412c7",
          "message": "Batch output in step-invariant operator (#132)\n\n* Batch output in step-invariant operator\r\n\r\nThe step invariant operator copies an input vector and replicates it\r\nover the entire query range. In the current implementation, the output\r\nof this operator is not chunked into batches, but rather expanded completely\r\ninto the total number of steps.\r\n\r\nThis leads to the binary operator not working properly since one operand\r\ncould produce 10 steps (the current batch size), and the other one\r\ncould produce steps for the entire range.\r\n\r\nThis commit changes the step invariant operator to produce output vectors\r\nthat are capped to the stepsBatch constant which is used in every other operator.\r\n\r\n* Update execution/step_invariant/step_invariant.go\r\n\r\nCo-authored-by: Saswata Mukherjee <saswataminsta@yahoo.com>\r\n\r\n* Import errors package\r\n\r\n* Update benchmarks\r\n\r\nCo-authored-by: Saswata Mukherjee <saswataminsta@yahoo.com>",
          "timestamp": "2022-12-05T07:19:37+01:00",
          "tree_id": "b5de48516acec867ad281787b960f2c21347ba5e",
          "url": "https://github.com/thanos-community/promql-engine/commit/fa1f96080ee7c732ad8f290db83ad7559aa412c7"
        },
        "date": 1670221349595,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26147171,
            "unit": "ns/op\t38721296 B/op\t  131257 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26668787,
            "unit": "ns/op\t38745445 B/op\t  131247 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26623713,
            "unit": "ns/op\t38735280 B/op\t  131236 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26554174,
            "unit": "ns/op\t38782339 B/op\t  131282 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26684193,
            "unit": "ns/op\t38714094 B/op\t  131221 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12308699,
            "unit": "ns/op\t12191491 B/op\t  125873 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12335365,
            "unit": "ns/op\t12186925 B/op\t  125890 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12383200,
            "unit": "ns/op\t12182101 B/op\t  125865 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12381617,
            "unit": "ns/op\t12164756 B/op\t  125863 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12372299,
            "unit": "ns/op\t12217954 B/op\t  125880 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21209148,
            "unit": "ns/op\t24266771 B/op\t  211646 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21079512,
            "unit": "ns/op\t24251828 B/op\t  211643 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21058087,
            "unit": "ns/op\t24233309 B/op\t  211632 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21220089,
            "unit": "ns/op\t24241082 B/op\t  211636 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21178980,
            "unit": "ns/op\t24258783 B/op\t  211636 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16722163,
            "unit": "ns/op\t11992384 B/op\t  135246 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16736903,
            "unit": "ns/op\t11966166 B/op\t  135205 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16710857,
            "unit": "ns/op\t12034631 B/op\t  135308 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16667222,
            "unit": "ns/op\t11984842 B/op\t  135181 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16734990,
            "unit": "ns/op\t11991033 B/op\t  135254 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16800945,
            "unit": "ns/op\t11986230 B/op\t  133322 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16756383,
            "unit": "ns/op\t12001668 B/op\t  133311 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16729406,
            "unit": "ns/op\t11971500 B/op\t  133313 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16749083,
            "unit": "ns/op\t11971370 B/op\t  133317 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16701728,
            "unit": "ns/op\t11988942 B/op\t  133334 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33420086,
            "unit": "ns/op\t41541626 B/op\t  155060 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33765841,
            "unit": "ns/op\t41573829 B/op\t  155100 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33971592,
            "unit": "ns/op\t41554099 B/op\t  155060 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33650477,
            "unit": "ns/op\t41620752 B/op\t  155130 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33875974,
            "unit": "ns/op\t41558650 B/op\t  155082 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23936954,
            "unit": "ns/op\t14637074 B/op\t  149655 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24224303,
            "unit": "ns/op\t14653494 B/op\t  149638 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24117674,
            "unit": "ns/op\t14709547 B/op\t  149690 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24085124,
            "unit": "ns/op\t14691829 B/op\t  149667 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24115889,
            "unit": "ns/op\t14669496 B/op\t  149629 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33031334,
            "unit": "ns/op\t26100634 B/op\t  235002 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33000109,
            "unit": "ns/op\t26165515 B/op\t  235065 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32934236,
            "unit": "ns/op\t26059982 B/op\t  234980 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32982317,
            "unit": "ns/op\t26066720 B/op\t  234987 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32976241,
            "unit": "ns/op\t26053310 B/op\t  234981 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30530555,
            "unit": "ns/op\t44468452 B/op\t  619133 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30550555,
            "unit": "ns/op\t44430841 B/op\t  619141 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30522203,
            "unit": "ns/op\t44425650 B/op\t  619126 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30569973,
            "unit": "ns/op\t44463872 B/op\t  619137 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30461462,
            "unit": "ns/op\t44399751 B/op\t  619123 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21168911,
            "unit": "ns/op\t22915828 B/op\t  108265 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21214318,
            "unit": "ns/op\t22922888 B/op\t  108306 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21108698,
            "unit": "ns/op\t22915748 B/op\t  108287 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21138777,
            "unit": "ns/op\t22929478 B/op\t  108304 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21162442,
            "unit": "ns/op\t22923617 B/op\t  108277 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 47046047,
            "unit": "ns/op\t57848214 B/op\t  202314 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46657367,
            "unit": "ns/op\t57828888 B/op\t  202262 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46755816,
            "unit": "ns/op\t57832818 B/op\t  202279 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46499412,
            "unit": "ns/op\t57835211 B/op\t  202291 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46744611,
            "unit": "ns/op\t57832789 B/op\t  202308 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36957809,
            "unit": "ns/op\t42279416 B/op\t  135562 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36824166,
            "unit": "ns/op\t42378409 B/op\t  135627 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37047260,
            "unit": "ns/op\t42280246 B/op\t  135576 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36612187,
            "unit": "ns/op\t42296191 B/op\t  135599 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36659027,
            "unit": "ns/op\t42343095 B/op\t  135617 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28786526,
            "unit": "ns/op\t40018362 B/op\t  143520 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28967245,
            "unit": "ns/op\t39987128 B/op\t  143487 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28827083,
            "unit": "ns/op\t40023343 B/op\t  143524 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28991232,
            "unit": "ns/op\t39986311 B/op\t  143476 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28766361,
            "unit": "ns/op\t40011983 B/op\t  143505 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37631461,
            "unit": "ns/op\t41903482 B/op\t  132557 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36594572,
            "unit": "ns/op\t41963335 B/op\t  132600 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36444522,
            "unit": "ns/op\t41921966 B/op\t  132571 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37143384,
            "unit": "ns/op\t41941454 B/op\t  132583 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37479764,
            "unit": "ns/op\t41858349 B/op\t  132545 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26787258,
            "unit": "ns/op\t36922274 B/op\t  102349 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26626733,
            "unit": "ns/op\t37013770 B/op\t  102420 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26787279,
            "unit": "ns/op\t36971548 B/op\t  102368 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26813519,
            "unit": "ns/op\t37017378 B/op\t  102410 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26618079,
            "unit": "ns/op\t36969524 B/op\t  102367 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22560205,
            "unit": "ns/op\t33140206 B/op\t   74749 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21906461,
            "unit": "ns/op\t33141477 B/op\t   74749 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22086314,
            "unit": "ns/op\t33139224 B/op\t   74749 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22323191,
            "unit": "ns/op\t33138838 B/op\t   74748 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22139119,
            "unit": "ns/op\t33140743 B/op\t   74749 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21086966,
            "unit": "ns/op\t32965012 B/op\t   68751 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21336110,
            "unit": "ns/op\t32960618 B/op\t   68751 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21255013,
            "unit": "ns/op\t32964037 B/op\t   68752 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 22179786,
            "unit": "ns/op\t32961121 B/op\t   68750 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20972290,
            "unit": "ns/op\t32960474 B/op\t   68750 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44486243,
            "unit": "ns/op\t39101982 B/op\t  135004 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44742646,
            "unit": "ns/op\t39106731 B/op\t  135065 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43956348,
            "unit": "ns/op\t39151250 B/op\t  135105 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43760512,
            "unit": "ns/op\t39147123 B/op\t  135100 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44447074,
            "unit": "ns/op\t39116375 B/op\t  135051 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41895421,
            "unit": "ns/op\t39082912 B/op\t  134635 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41874518,
            "unit": "ns/op\t39076192 B/op\t  134622 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41812660,
            "unit": "ns/op\t39050288 B/op\t  134594 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42438521,
            "unit": "ns/op\t39062506 B/op\t  134617 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41464552,
            "unit": "ns/op\t39084539 B/op\t  134627 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53991730,
            "unit": "ns/op\t42595078 B/op\t  139494 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55063149,
            "unit": "ns/op\t42606517 B/op\t  139497 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55375429,
            "unit": "ns/op\t42497735 B/op\t  139420 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 56073623,
            "unit": "ns/op\t42550401 B/op\t  139463 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55510081,
            "unit": "ns/op\t42491963 B/op\t  139422 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46903829,
            "unit": "ns/op\t41568657 B/op\t  156099 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47018685,
            "unit": "ns/op\t41541488 B/op\t  156038 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47013084,
            "unit": "ns/op\t41549443 B/op\t  156066 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46898267,
            "unit": "ns/op\t41558511 B/op\t  156101 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46881810,
            "unit": "ns/op\t41588996 B/op\t  156103 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47498841,
            "unit": "ns/op\t41519708 B/op\t  156034 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47178971,
            "unit": "ns/op\t41587520 B/op\t  156104 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47370522,
            "unit": "ns/op\t41531839 B/op\t  156030 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47220164,
            "unit": "ns/op\t41520432 B/op\t  156033 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46834076,
            "unit": "ns/op\t41564256 B/op\t  156065 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129849621,
            "unit": "ns/op\t60486609 B/op\t  704586 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130127614,
            "unit": "ns/op\t60447439 B/op\t  704605 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130152355,
            "unit": "ns/op\t60382749 B/op\t  704541 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129494563,
            "unit": "ns/op\t60533428 B/op\t  704583 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130700344,
            "unit": "ns/op\t60475786 B/op\t  704593 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6f0bb4bec09c7df2849cd3735a2ff81ce84422d6",
          "message": "engine: make histogram_quantile work on data without le label (#129)\n\n* engine: add histogram_quantile test case with malformed data\r\n\r\n* histogram: check for index\r\n\r\n* engine: add more histogram_quantile tests\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>\r\n\r\n* function: rework histogram fix\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>\r\n\r\n* function: fix according to suggestions\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>",
          "timestamp": "2022-12-05T11:10:37+02:00",
          "tree_id": "0b7980ffea051b2973671db227fec3894720ae30",
          "url": "https://github.com/thanos-community/promql-engine/commit/6f0bb4bec09c7df2849cd3735a2ff81ce84422d6"
        },
        "date": 1670231606106,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25505449,
            "unit": "ns/op\t38793602 B/op\t  131646 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26371760,
            "unit": "ns/op\t38774372 B/op\t  131616 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26614195,
            "unit": "ns/op\t38772550 B/op\t  131607 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26576730,
            "unit": "ns/op\t38740783 B/op\t  131588 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26676751,
            "unit": "ns/op\t38752843 B/op\t  131604 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12063228,
            "unit": "ns/op\t12159380 B/op\t  126191 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12064066,
            "unit": "ns/op\t12235377 B/op\t  126238 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11991418,
            "unit": "ns/op\t12227341 B/op\t  126242 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12029616,
            "unit": "ns/op\t12233199 B/op\t  126252 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12035303,
            "unit": "ns/op\t12227167 B/op\t  126223 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21313953,
            "unit": "ns/op\t24274963 B/op\t  211996 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21148850,
            "unit": "ns/op\t24281133 B/op\t  212002 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21281900,
            "unit": "ns/op\t24267563 B/op\t  211991 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21220650,
            "unit": "ns/op\t24278911 B/op\t  211994 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21212278,
            "unit": "ns/op\t24265079 B/op\t  212008 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16487149,
            "unit": "ns/op\t12086321 B/op\t  133981 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16419620,
            "unit": "ns/op\t12092471 B/op\t  133987 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16523270,
            "unit": "ns/op\t12061261 B/op\t  134021 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16381909,
            "unit": "ns/op\t12017937 B/op\t  133933 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16423015,
            "unit": "ns/op\t12017419 B/op\t  133937 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16464760,
            "unit": "ns/op\t12079795 B/op\t  134633 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16518335,
            "unit": "ns/op\t12102323 B/op\t  134672 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16499346,
            "unit": "ns/op\t12104765 B/op\t  134683 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16465999,
            "unit": "ns/op\t12099995 B/op\t  134657 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16490963,
            "unit": "ns/op\t12040641 B/op\t  134631 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33208464,
            "unit": "ns/op\t41583602 B/op\t  155470 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33343502,
            "unit": "ns/op\t41611813 B/op\t  155480 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33572076,
            "unit": "ns/op\t41571165 B/op\t  155429 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33346055,
            "unit": "ns/op\t41554960 B/op\t  155446 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33194834,
            "unit": "ns/op\t41587175 B/op\t  155445 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23650335,
            "unit": "ns/op\t14709119 B/op\t  150020 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23772213,
            "unit": "ns/op\t14664451 B/op\t  149996 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23840098,
            "unit": "ns/op\t14695800 B/op\t  150039 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23708068,
            "unit": "ns/op\t14692758 B/op\t  150022 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23638473,
            "unit": "ns/op\t14693094 B/op\t  150026 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32648280,
            "unit": "ns/op\t26255508 B/op\t  235500 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32922325,
            "unit": "ns/op\t26228625 B/op\t  235455 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32689532,
            "unit": "ns/op\t26170672 B/op\t  235421 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32713638,
            "unit": "ns/op\t26168272 B/op\t  235427 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32811547,
            "unit": "ns/op\t26133595 B/op\t  235396 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30637855,
            "unit": "ns/op\t44518820 B/op\t  619849 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30488160,
            "unit": "ns/op\t44560609 B/op\t  619879 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30323151,
            "unit": "ns/op\t44507100 B/op\t  619850 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30371034,
            "unit": "ns/op\t44415658 B/op\t  619811 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30268336,
            "unit": "ns/op\t44455298 B/op\t  619845 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21144934,
            "unit": "ns/op\t22926414 B/op\t  108578 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21234013,
            "unit": "ns/op\t22932825 B/op\t  108602 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21246886,
            "unit": "ns/op\t22941270 B/op\t  108614 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21159926,
            "unit": "ns/op\t22926841 B/op\t  108587 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21030932,
            "unit": "ns/op\t22932102 B/op\t  108576 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46997706,
            "unit": "ns/op\t57857687 B/op\t  202757 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46679860,
            "unit": "ns/op\t57816545 B/op\t  202734 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46313899,
            "unit": "ns/op\t57864945 B/op\t  202772 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46672831,
            "unit": "ns/op\t57840132 B/op\t  202750 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46855062,
            "unit": "ns/op\t57853224 B/op\t  202749 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36984876,
            "unit": "ns/op\t42302108 B/op\t  135922 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37191662,
            "unit": "ns/op\t42346763 B/op\t  135966 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36754923,
            "unit": "ns/op\t42428121 B/op\t  136022 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37131169,
            "unit": "ns/op\t42387857 B/op\t  135982 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36723468,
            "unit": "ns/op\t42439983 B/op\t  136018 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28762760,
            "unit": "ns/op\t40015045 B/op\t  143870 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28854799,
            "unit": "ns/op\t40025168 B/op\t  143869 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28491186,
            "unit": "ns/op\t40057263 B/op\t  143886 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28388901,
            "unit": "ns/op\t40016774 B/op\t  143870 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28629748,
            "unit": "ns/op\t40016276 B/op\t  143861 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37494758,
            "unit": "ns/op\t41925847 B/op\t  132922 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37200708,
            "unit": "ns/op\t41981865 B/op\t  132951 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36916221,
            "unit": "ns/op\t41974906 B/op\t  132953 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36957022,
            "unit": "ns/op\t41975581 B/op\t  132951 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37185420,
            "unit": "ns/op\t41898356 B/op\t  132914 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26623595,
            "unit": "ns/op\t37055198 B/op\t  102666 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26616654,
            "unit": "ns/op\t37002941 B/op\t  102617 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26916437,
            "unit": "ns/op\t36991839 B/op\t  102624 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26404113,
            "unit": "ns/op\t36982665 B/op\t  102613 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26713440,
            "unit": "ns/op\t36985632 B/op\t  102623 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22026609,
            "unit": "ns/op\t33146990 B/op\t   74870 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21722307,
            "unit": "ns/op\t33143969 B/op\t   74868 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21263125,
            "unit": "ns/op\t33146445 B/op\t   74870 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21907854,
            "unit": "ns/op\t33143291 B/op\t   74869 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22174750,
            "unit": "ns/op\t33143164 B/op\t   74868 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20970716,
            "unit": "ns/op\t32973284 B/op\t   68873 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21038178,
            "unit": "ns/op\t32970963 B/op\t   68871 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21190704,
            "unit": "ns/op\t32965166 B/op\t   68869 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21035547,
            "unit": "ns/op\t32965925 B/op\t   68869 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21015924,
            "unit": "ns/op\t32972586 B/op\t   68871 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44921478,
            "unit": "ns/op\t39095076 B/op\t  135359 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45151041,
            "unit": "ns/op\t39135648 B/op\t  135404 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44467528,
            "unit": "ns/op\t39162248 B/op\t  135449 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45308322,
            "unit": "ns/op\t39108313 B/op\t  135365 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45110513,
            "unit": "ns/op\t39140890 B/op\t  135397 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41789474,
            "unit": "ns/op\t39090972 B/op\t  134970 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41730161,
            "unit": "ns/op\t39091395 B/op\t  134989 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41087232,
            "unit": "ns/op\t39127282 B/op\t  134999 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41954244,
            "unit": "ns/op\t39087341 B/op\t  134975 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41316880,
            "unit": "ns/op\t39107284 B/op\t  134998 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55354143,
            "unit": "ns/op\t42463000 B/op\t  139747 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55070596,
            "unit": "ns/op\t42489915 B/op\t  139782 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54789825,
            "unit": "ns/op\t42630951 B/op\t  139862 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54968726,
            "unit": "ns/op\t42545337 B/op\t  139832 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55979898,
            "unit": "ns/op\t42452905 B/op\t  139753 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47214051,
            "unit": "ns/op\t41596958 B/op\t  156447 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46766889,
            "unit": "ns/op\t41572454 B/op\t  156423 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46216835,
            "unit": "ns/op\t41603706 B/op\t  156458 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46554885,
            "unit": "ns/op\t41612718 B/op\t  156476 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46633761,
            "unit": "ns/op\t41601866 B/op\t  156464 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47066034,
            "unit": "ns/op\t41597399 B/op\t  156470 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46902350,
            "unit": "ns/op\t41563983 B/op\t  156417 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47085310,
            "unit": "ns/op\t41597315 B/op\t  156464 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47149282,
            "unit": "ns/op\t41610554 B/op\t  156486 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47166849,
            "unit": "ns/op\t41566654 B/op\t  156433 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129985322,
            "unit": "ns/op\t60555851 B/op\t  715186 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128956716,
            "unit": "ns/op\t60456825 B/op\t  715120 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130538191,
            "unit": "ns/op\t60528995 B/op\t  715172 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131195957,
            "unit": "ns/op\t60634800 B/op\t  715177 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130245025,
            "unit": "ns/op\t60625504 B/op\t  715169 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "effd058456da0781a00258af5921cd617600b63b",
          "message": "engine: make histogram_quantile resilient against duplicate output IDs (#135)\n\no.seriesBuckets[i] can be nil or empty if multiple input series map to\r\nthe same output series IDs. Skip them over in such case.",
          "timestamp": "2022-12-07T12:17:42+02:00",
          "tree_id": "3fa8cd0df2f0afd6a669192a290122069b984208",
          "url": "https://github.com/thanos-community/promql-engine/commit/effd058456da0781a00258af5921cd617600b63b"
        },
        "date": 1670408430136,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26207068,
            "unit": "ns/op\t38741632 B/op\t  131555 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26289968,
            "unit": "ns/op\t38767228 B/op\t  131579 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26409757,
            "unit": "ns/op\t38766369 B/op\t  131550 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26478230,
            "unit": "ns/op\t38803786 B/op\t  131578 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26708972,
            "unit": "ns/op\t38777588 B/op\t  131580 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12028641,
            "unit": "ns/op\t12194847 B/op\t  126166 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12038192,
            "unit": "ns/op\t12236160 B/op\t  126201 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12097920,
            "unit": "ns/op\t12226859 B/op\t  126190 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12064840,
            "unit": "ns/op\t12202282 B/op\t  126170 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12061496,
            "unit": "ns/op\t12173358 B/op\t  126182 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21312399,
            "unit": "ns/op\t24262922 B/op\t  211941 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21114986,
            "unit": "ns/op\t24259362 B/op\t  211936 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21148129,
            "unit": "ns/op\t24276456 B/op\t  211946 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21170314,
            "unit": "ns/op\t24252681 B/op\t  211945 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21251229,
            "unit": "ns/op\t24268854 B/op\t  211948 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16487132,
            "unit": "ns/op\t12089475 B/op\t  134178 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16504528,
            "unit": "ns/op\t12032010 B/op\t  134175 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16466544,
            "unit": "ns/op\t12023781 B/op\t  134133 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16467404,
            "unit": "ns/op\t12000826 B/op\t  134132 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16506039,
            "unit": "ns/op\t12109008 B/op\t  134169 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16566439,
            "unit": "ns/op\t12094337 B/op\t  134516 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16478495,
            "unit": "ns/op\t12050413 B/op\t  134520 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16562807,
            "unit": "ns/op\t12102973 B/op\t  134585 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16463363,
            "unit": "ns/op\t12082367 B/op\t  134552 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16539696,
            "unit": "ns/op\t12071054 B/op\t  134510 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33239908,
            "unit": "ns/op\t41552520 B/op\t  155360 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33273342,
            "unit": "ns/op\t41581936 B/op\t  155399 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33208128,
            "unit": "ns/op\t41598854 B/op\t  155417 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33473252,
            "unit": "ns/op\t41606026 B/op\t  155413 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33245484,
            "unit": "ns/op\t41574687 B/op\t  155385 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23796405,
            "unit": "ns/op\t14616538 B/op\t  149905 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23734489,
            "unit": "ns/op\t14742369 B/op\t  150005 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23729308,
            "unit": "ns/op\t14686415 B/op\t  150000 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23652367,
            "unit": "ns/op\t14745951 B/op\t  150031 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23918905,
            "unit": "ns/op\t14651684 B/op\t  149913 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33004186,
            "unit": "ns/op\t26110246 B/op\t  235316 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32680478,
            "unit": "ns/op\t26167370 B/op\t  235372 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32824195,
            "unit": "ns/op\t26127177 B/op\t  235329 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32826179,
            "unit": "ns/op\t26165538 B/op\t  235362 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32775263,
            "unit": "ns/op\t26148111 B/op\t  235352 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30764945,
            "unit": "ns/op\t44490908 B/op\t  619762 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30273259,
            "unit": "ns/op\t44458844 B/op\t  619718 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30541489,
            "unit": "ns/op\t44566689 B/op\t  619768 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30211194,
            "unit": "ns/op\t44515624 B/op\t  619748 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30550848,
            "unit": "ns/op\t44520847 B/op\t  619769 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21224596,
            "unit": "ns/op\t22934433 B/op\t  108482 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21195755,
            "unit": "ns/op\t22937337 B/op\t  108470 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21142888,
            "unit": "ns/op\t22910939 B/op\t  108448 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21174853,
            "unit": "ns/op\t22914544 B/op\t  108432 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21231945,
            "unit": "ns/op\t22919279 B/op\t  108459 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46735792,
            "unit": "ns/op\t57798725 B/op\t  202696 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 47015662,
            "unit": "ns/op\t57842059 B/op\t  202732 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46919297,
            "unit": "ns/op\t57840827 B/op\t  202749 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46783347,
            "unit": "ns/op\t57867439 B/op\t  202774 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46863509,
            "unit": "ns/op\t57855355 B/op\t  202757 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37187641,
            "unit": "ns/op\t42328528 B/op\t  135899 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36931996,
            "unit": "ns/op\t42287846 B/op\t  135894 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37007751,
            "unit": "ns/op\t42295891 B/op\t  135883 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36958797,
            "unit": "ns/op\t42294442 B/op\t  135869 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36400206,
            "unit": "ns/op\t42360117 B/op\t  135914 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 29209693,
            "unit": "ns/op\t40012875 B/op\t  143812 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28641354,
            "unit": "ns/op\t40046605 B/op\t  143827 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28595155,
            "unit": "ns/op\t40026217 B/op\t  143816 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28434059,
            "unit": "ns/op\t40076417 B/op\t  143874 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28553933,
            "unit": "ns/op\t40072371 B/op\t  143856 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37148942,
            "unit": "ns/op\t41896452 B/op\t  132855 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36917969,
            "unit": "ns/op\t41865903 B/op\t  132840 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36764516,
            "unit": "ns/op\t41948491 B/op\t  132888 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36772563,
            "unit": "ns/op\t41988394 B/op\t  132899 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37467628,
            "unit": "ns/op\t41969123 B/op\t  132886 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26998195,
            "unit": "ns/op\t36986003 B/op\t  102593 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26602617,
            "unit": "ns/op\t37067691 B/op\t  102662 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26783839,
            "unit": "ns/op\t37001733 B/op\t  102598 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26613035,
            "unit": "ns/op\t37027339 B/op\t  102626 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26666848,
            "unit": "ns/op\t36985375 B/op\t  102583 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22069352,
            "unit": "ns/op\t33146775 B/op\t   74852 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21765340,
            "unit": "ns/op\t33146016 B/op\t   74851 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21057209,
            "unit": "ns/op\t33143360 B/op\t   74850 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21119237,
            "unit": "ns/op\t33142100 B/op\t   74850 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21364617,
            "unit": "ns/op\t33142852 B/op\t   74851 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21218682,
            "unit": "ns/op\t32965748 B/op\t   68853 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21250619,
            "unit": "ns/op\t32967925 B/op\t   68853 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20936411,
            "unit": "ns/op\t32968078 B/op\t   68853 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21413352,
            "unit": "ns/op\t32970670 B/op\t   68854 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20684473,
            "unit": "ns/op\t32965104 B/op\t   68854 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45154246,
            "unit": "ns/op\t39125408 B/op\t  135358 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45006896,
            "unit": "ns/op\t39139920 B/op\t  135364 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44839964,
            "unit": "ns/op\t39114025 B/op\t  135318 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45043103,
            "unit": "ns/op\t39128294 B/op\t  135334 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44555069,
            "unit": "ns/op\t39148244 B/op\t  135365 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41705151,
            "unit": "ns/op\t39105752 B/op\t  134946 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41823558,
            "unit": "ns/op\t39089925 B/op\t  134932 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42072392,
            "unit": "ns/op\t39114885 B/op\t  134962 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41980460,
            "unit": "ns/op\t39118386 B/op\t  134951 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41951984,
            "unit": "ns/op\t39084189 B/op\t  134936 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54918657,
            "unit": "ns/op\t42525579 B/op\t  139739 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55743226,
            "unit": "ns/op\t42438019 B/op\t  139694 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55309270,
            "unit": "ns/op\t42554419 B/op\t  139758 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55323682,
            "unit": "ns/op\t42501780 B/op\t  139721 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55698660,
            "unit": "ns/op\t42524486 B/op\t  139747 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47512935,
            "unit": "ns/op\t41576050 B/op\t  156385 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47261299,
            "unit": "ns/op\t41605574 B/op\t  156448 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47122837,
            "unit": "ns/op\t41600046 B/op\t  156412 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47218456,
            "unit": "ns/op\t41564496 B/op\t  156400 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47634135,
            "unit": "ns/op\t41527304 B/op\t  156312 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47356495,
            "unit": "ns/op\t41582085 B/op\t  156407 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47526442,
            "unit": "ns/op\t41555160 B/op\t  156378 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47824547,
            "unit": "ns/op\t41555565 B/op\t  156376 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47050471,
            "unit": "ns/op\t41575966 B/op\t  156406 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47158707,
            "unit": "ns/op\t41579153 B/op\t  156391 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129924931,
            "unit": "ns/op\t60493836 B/op\t  715324 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128241692,
            "unit": "ns/op\t60596224 B/op\t  715323 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130616111,
            "unit": "ns/op\t60478217 B/op\t  715269 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 128936018,
            "unit": "ns/op\t60451789 B/op\t  715278 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 132053608,
            "unit": "ns/op\t60434378 B/op\t  715271 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "05248ae34b960b90766c2aaa4611bc5882da58aa",
          "message": "Fix topk/bottomk on empty inputs (#136)\n\nTopK and bottomK can take in empty input vectors in certain cases,\r\nwhich can lead to validation failing in the kHashAggregate operator.\r\n\r\nThis commit changes the validation to check if the input vector batch has\r\nless steps than the scalar vector batch.",
          "timestamp": "2022-12-08T08:12:56+01:00",
          "tree_id": "486c22ba1601620dcdb9a3997adeca5a75d80c1c",
          "url": "https://github.com/thanos-community/promql-engine/commit/05248ae34b960b90766c2aaa4611bc5882da58aa"
        },
        "date": 1670483746004,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25765738,
            "unit": "ns/op\t38786727 B/op\t  131761 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26363708,
            "unit": "ns/op\t38792064 B/op\t  131728 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26507238,
            "unit": "ns/op\t38785019 B/op\t  131733 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26578415,
            "unit": "ns/op\t38785030 B/op\t  131730 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26454172,
            "unit": "ns/op\t38827342 B/op\t  131741 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12085777,
            "unit": "ns/op\t12250477 B/op\t  126348 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12038746,
            "unit": "ns/op\t12236086 B/op\t  126354 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12029317,
            "unit": "ns/op\t12193947 B/op\t  126335 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12036157,
            "unit": "ns/op\t12242679 B/op\t  126359 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12027401,
            "unit": "ns/op\t12191429 B/op\t  126322 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21124242,
            "unit": "ns/op\t24255515 B/op\t  212079 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21015700,
            "unit": "ns/op\t24267145 B/op\t  212114 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21089308,
            "unit": "ns/op\t24270980 B/op\t  212099 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21160088,
            "unit": "ns/op\t24269870 B/op\t  212105 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21141961,
            "unit": "ns/op\t24246814 B/op\t  212096 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16465051,
            "unit": "ns/op\t12080193 B/op\t  135186 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16478207,
            "unit": "ns/op\t12124356 B/op\t  135203 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16452999,
            "unit": "ns/op\t12055735 B/op\t  135154 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16529784,
            "unit": "ns/op\t12086001 B/op\t  135217 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16528289,
            "unit": "ns/op\t12126783 B/op\t  135258 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16411752,
            "unit": "ns/op\t12081656 B/op\t  133671 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16479245,
            "unit": "ns/op\t12062941 B/op\t  133645 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16457663,
            "unit": "ns/op\t12072058 B/op\t  133674 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16395794,
            "unit": "ns/op\t12059974 B/op\t  133638 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16375359,
            "unit": "ns/op\t12033106 B/op\t  133626 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33013788,
            "unit": "ns/op\t41596608 B/op\t  155568 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33140371,
            "unit": "ns/op\t41577090 B/op\t  155547 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33103269,
            "unit": "ns/op\t41589094 B/op\t  155540 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33283511,
            "unit": "ns/op\t41633847 B/op\t  155613 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33291885,
            "unit": "ns/op\t41601494 B/op\t  155555 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23788087,
            "unit": "ns/op\t14664954 B/op\t  150073 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23874232,
            "unit": "ns/op\t14711027 B/op\t  150096 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23744643,
            "unit": "ns/op\t14746625 B/op\t  150185 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23800043,
            "unit": "ns/op\t14758284 B/op\t  150180 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23639251,
            "unit": "ns/op\t14720366 B/op\t  150118 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32518778,
            "unit": "ns/op\t26206278 B/op\t  235541 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32610773,
            "unit": "ns/op\t26175760 B/op\t  235533 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32522407,
            "unit": "ns/op\t26097009 B/op\t  235462 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32583228,
            "unit": "ns/op\t26148545 B/op\t  235494 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32666454,
            "unit": "ns/op\t26141251 B/op\t  235494 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30443968,
            "unit": "ns/op\t44550748 B/op\t  620078 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30337682,
            "unit": "ns/op\t44525258 B/op\t  620078 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30195181,
            "unit": "ns/op\t44479270 B/op\t  620053 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30364398,
            "unit": "ns/op\t44463306 B/op\t  620053 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30374803,
            "unit": "ns/op\t44576213 B/op\t  620101 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21190082,
            "unit": "ns/op\t22922630 B/op\t  108541 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21184826,
            "unit": "ns/op\t22941188 B/op\t  108565 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21243282,
            "unit": "ns/op\t22924514 B/op\t  108549 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21322889,
            "unit": "ns/op\t22940015 B/op\t  108581 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 21167477,
            "unit": "ns/op\t22943120 B/op\t  108578 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 47051302,
            "unit": "ns/op\t57886931 B/op\t  202785 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46973168,
            "unit": "ns/op\t57843349 B/op\t  202736 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46855325,
            "unit": "ns/op\t57944214 B/op\t  202867 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46389032,
            "unit": "ns/op\t57871773 B/op\t  202760 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 46380657,
            "unit": "ns/op\t57837049 B/op\t  202705 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37259383,
            "unit": "ns/op\t42329783 B/op\t  136049 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36994162,
            "unit": "ns/op\t42306747 B/op\t  136051 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36875534,
            "unit": "ns/op\t42395076 B/op\t  136092 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36511948,
            "unit": "ns/op\t42453076 B/op\t  136131 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37120459,
            "unit": "ns/op\t42377663 B/op\t  136086 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28399904,
            "unit": "ns/op\t40022054 B/op\t  143972 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28456091,
            "unit": "ns/op\t40053248 B/op\t  143996 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28654660,
            "unit": "ns/op\t40018904 B/op\t  143991 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28400036,
            "unit": "ns/op\t40044817 B/op\t  143980 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28338113,
            "unit": "ns/op\t40076157 B/op\t  143996 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37525517,
            "unit": "ns/op\t41825130 B/op\t  132981 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36874391,
            "unit": "ns/op\t41989104 B/op\t  133063 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36767288,
            "unit": "ns/op\t41946001 B/op\t  133044 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36852474,
            "unit": "ns/op\t41934319 B/op\t  133031 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36826230,
            "unit": "ns/op\t41941237 B/op\t  133034 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26751460,
            "unit": "ns/op\t37050239 B/op\t  102737 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26807028,
            "unit": "ns/op\t37006383 B/op\t  102707 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26601348,
            "unit": "ns/op\t37006570 B/op\t  102708 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26264152,
            "unit": "ns/op\t37042727 B/op\t  102734 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26449111,
            "unit": "ns/op\t37023904 B/op\t  102716 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22394272,
            "unit": "ns/op\t33148460 B/op\t   74904 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21943925,
            "unit": "ns/op\t33149118 B/op\t   74904 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21458510,
            "unit": "ns/op\t33146781 B/op\t   74905 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22455399,
            "unit": "ns/op\t33145674 B/op\t   74905 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21583495,
            "unit": "ns/op\t33148373 B/op\t   74904 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21300990,
            "unit": "ns/op\t32964723 B/op\t   68905 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20924601,
            "unit": "ns/op\t32973300 B/op\t   68907 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21193963,
            "unit": "ns/op\t32973476 B/op\t   68907 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20764610,
            "unit": "ns/op\t32971352 B/op\t   68908 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20946842,
            "unit": "ns/op\t32966928 B/op\t   68905 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44736074,
            "unit": "ns/op\t39153446 B/op\t  135525 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43828414,
            "unit": "ns/op\t39174075 B/op\t  135564 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44684106,
            "unit": "ns/op\t39155720 B/op\t  135523 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44799585,
            "unit": "ns/op\t39145815 B/op\t  135475 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44368036,
            "unit": "ns/op\t39149365 B/op\t  135515 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42581290,
            "unit": "ns/op\t39075483 B/op\t  135061 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41760925,
            "unit": "ns/op\t39101027 B/op\t  135097 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41982112,
            "unit": "ns/op\t39104173 B/op\t  135118 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42462203,
            "unit": "ns/op\t39094850 B/op\t  135086 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42090848,
            "unit": "ns/op\t39080002 B/op\t  135070 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55949264,
            "unit": "ns/op\t42428632 B/op\t  139835 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54975412,
            "unit": "ns/op\t42457198 B/op\t  139872 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54581090,
            "unit": "ns/op\t42491810 B/op\t  139880 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54467215,
            "unit": "ns/op\t42571512 B/op\t  139931 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54478981,
            "unit": "ns/op\t42423285 B/op\t  139828 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46824822,
            "unit": "ns/op\t41636936 B/op\t  156600 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47208452,
            "unit": "ns/op\t41574664 B/op\t  156517 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46955128,
            "unit": "ns/op\t41616232 B/op\t  156581 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47569144,
            "unit": "ns/op\t41543013 B/op\t  156491 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47550646,
            "unit": "ns/op\t41556454 B/op\t  156496 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47218031,
            "unit": "ns/op\t41556424 B/op\t  156505 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47147537,
            "unit": "ns/op\t41572750 B/op\t  156538 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46770496,
            "unit": "ns/op\t41628092 B/op\t  156596 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46588099,
            "unit": "ns/op\t41565603 B/op\t  156537 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46818694,
            "unit": "ns/op\t41593698 B/op\t  156564 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129385316,
            "unit": "ns/op\t60423321 B/op\t  714853 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129723809,
            "unit": "ns/op\t60618476 B/op\t  714919 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130698454,
            "unit": "ns/op\t60714702 B/op\t  714975 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 131333024,
            "unit": "ns/op\t60537138 B/op\t  714940 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 127651365,
            "unit": "ns/op\t60577196 B/op\t  714909 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "csmarchbanks@gmail.com",
            "name": "Chris Marchbanks",
            "username": "csmarchbanks"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "654bfb6f34ffc49f5c62aa96154f368cf733b50d",
          "message": "Fetch Series in parallel for vector operators (#137)\n\n* Fetch Series in parallel for vector operators\r\n\r\nFixes a TODO in the code, just fetching the Series can be slow in a\r\nPromQL server I run and running the operations in parallel can halve the\r\nresponse time.\r\n\r\nSigned-off-by: Chris Marchbanks <csmarchbanks@gmail.com>\r\n\r\n* Update hints test to not worry about order\r\n\r\nSigned-off-by: Chris Marchbanks <csmarchbanks@gmail.com>",
          "timestamp": "2022-12-12T17:21:55+02:00",
          "tree_id": "fd638d42d1152d7a77ae324ebb7ef6f1a46b3cda",
          "url": "https://github.com/thanos-community/promql-engine/commit/654bfb6f34ffc49f5c62aa96154f368cf733b50d"
        },
        "date": 1670858687764,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26167535,
            "unit": "ns/op\t38736061 B/op\t  131226 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26140289,
            "unit": "ns/op\t38753789 B/op\t  131215 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26678179,
            "unit": "ns/op\t38781301 B/op\t  131236 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26366933,
            "unit": "ns/op\t38715363 B/op\t  131192 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26545733,
            "unit": "ns/op\t38806862 B/op\t  131247 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12097206,
            "unit": "ns/op\t12199936 B/op\t  125847 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12093088,
            "unit": "ns/op\t12195282 B/op\t  125841 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12012115,
            "unit": "ns/op\t12195310 B/op\t  125853 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12018910,
            "unit": "ns/op\t12222067 B/op\t  125856 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12020284,
            "unit": "ns/op\t12165523 B/op\t  125831 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21566674,
            "unit": "ns/op\t24256447 B/op\t  211607 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21260424,
            "unit": "ns/op\t24284854 B/op\t  211619 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20990659,
            "unit": "ns/op\t24259845 B/op\t  211596 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20951169,
            "unit": "ns/op\t24244939 B/op\t  211597 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20976088,
            "unit": "ns/op\t24257330 B/op\t  211584 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16502932,
            "unit": "ns/op\t12026548 B/op\t  132872 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16423126,
            "unit": "ns/op\t12032074 B/op\t  132866 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16369790,
            "unit": "ns/op\t12015457 B/op\t  132877 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16410261,
            "unit": "ns/op\t12093307 B/op\t  132969 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16424837,
            "unit": "ns/op\t12066511 B/op\t  132928 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16428882,
            "unit": "ns/op\t12053664 B/op\t  133880 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16391148,
            "unit": "ns/op\t12060366 B/op\t  133857 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16448499,
            "unit": "ns/op\t12044872 B/op\t  133877 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16420455,
            "unit": "ns/op\t12067977 B/op\t  133905 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16406720,
            "unit": "ns/op\t12109125 B/op\t  133942 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33123173,
            "unit": "ns/op\t41581115 B/op\t  155062 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33414892,
            "unit": "ns/op\t41567885 B/op\t  155060 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33476690,
            "unit": "ns/op\t41580496 B/op\t  155068 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33272431,
            "unit": "ns/op\t41556988 B/op\t  155014 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33459629,
            "unit": "ns/op\t41561759 B/op\t  155058 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23872786,
            "unit": "ns/op\t14711507 B/op\t  149653 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23880649,
            "unit": "ns/op\t14658837 B/op\t  149635 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24056231,
            "unit": "ns/op\t14694837 B/op\t  149613 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23979779,
            "unit": "ns/op\t14639344 B/op\t  149568 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23968356,
            "unit": "ns/op\t14641508 B/op\t  149593 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32658964,
            "unit": "ns/op\t26126899 B/op\t  235010 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32737632,
            "unit": "ns/op\t26030625 B/op\t  234918 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32814359,
            "unit": "ns/op\t26105136 B/op\t  234993 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32918583,
            "unit": "ns/op\t26056529 B/op\t  234947 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32616939,
            "unit": "ns/op\t26074386 B/op\t  234969 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30513679,
            "unit": "ns/op\t44430860 B/op\t  619097 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30388865,
            "unit": "ns/op\t44410268 B/op\t  619053 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30186788,
            "unit": "ns/op\t44486709 B/op\t  619087 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30401956,
            "unit": "ns/op\t44498438 B/op\t  619063 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30309061,
            "unit": "ns/op\t44416800 B/op\t  619024 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18905347,
            "unit": "ns/op\t22926923 B/op\t  108466 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18896559,
            "unit": "ns/op\t22917951 B/op\t  108464 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18956605,
            "unit": "ns/op\t22927718 B/op\t  108468 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18981125,
            "unit": "ns/op\t22927606 B/op\t  108479 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18859084,
            "unit": "ns/op\t22916494 B/op\t  108455 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44428005,
            "unit": "ns/op\t57858137 B/op\t  202448 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44761651,
            "unit": "ns/op\t57824456 B/op\t  202417 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43966687,
            "unit": "ns/op\t57771585 B/op\t  202348 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44403480,
            "unit": "ns/op\t57871375 B/op\t  202457 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44412886,
            "unit": "ns/op\t57819378 B/op\t  202406 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37329919,
            "unit": "ns/op\t42343313 B/op\t  135561 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36891085,
            "unit": "ns/op\t42316053 B/op\t  135553 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37193842,
            "unit": "ns/op\t42336424 B/op\t  135566 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37088055,
            "unit": "ns/op\t42389021 B/op\t  135595 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36787832,
            "unit": "ns/op\t42312882 B/op\t  135560 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28952088,
            "unit": "ns/op\t40003201 B/op\t  143472 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28744672,
            "unit": "ns/op\t40039610 B/op\t  143498 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28513886,
            "unit": "ns/op\t40017071 B/op\t  143464 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28205083,
            "unit": "ns/op\t40040056 B/op\t  143491 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28561099,
            "unit": "ns/op\t40039991 B/op\t  143497 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37487102,
            "unit": "ns/op\t41916095 B/op\t  132525 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37469626,
            "unit": "ns/op\t41775659 B/op\t  132455 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36821214,
            "unit": "ns/op\t41902242 B/op\t  132503 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36918314,
            "unit": "ns/op\t41828152 B/op\t  132491 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37173358,
            "unit": "ns/op\t41973653 B/op\t  132558 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26496149,
            "unit": "ns/op\t36974851 B/op\t  102349 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26672229,
            "unit": "ns/op\t36992507 B/op\t  102363 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26526500,
            "unit": "ns/op\t37017717 B/op\t  102396 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26636223,
            "unit": "ns/op\t37004882 B/op\t  102367 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26691569,
            "unit": "ns/op\t36986378 B/op\t  102367 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22116133,
            "unit": "ns/op\t33147231 B/op\t   74738 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21608319,
            "unit": "ns/op\t33145734 B/op\t   74737 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22064117,
            "unit": "ns/op\t33141281 B/op\t   74737 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21632879,
            "unit": "ns/op\t33142480 B/op\t   74736 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21433648,
            "unit": "ns/op\t33145333 B/op\t   74738 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20499393,
            "unit": "ns/op\t32966786 B/op\t   68738 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21138495,
            "unit": "ns/op\t32961883 B/op\t   68737 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20640412,
            "unit": "ns/op\t32966384 B/op\t   68738 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21160447,
            "unit": "ns/op\t32966170 B/op\t   68739 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21126389,
            "unit": "ns/op\t32962198 B/op\t   68738 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45084020,
            "unit": "ns/op\t39114371 B/op\t  135021 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43737802,
            "unit": "ns/op\t39142308 B/op\t  135029 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43967600,
            "unit": "ns/op\t39148646 B/op\t  135027 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 43693347,
            "unit": "ns/op\t39162005 B/op\t  135052 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44836920,
            "unit": "ns/op\t39123840 B/op\t  135035 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42010330,
            "unit": "ns/op\t39077012 B/op\t  134565 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41693580,
            "unit": "ns/op\t39115456 B/op\t  134629 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42308344,
            "unit": "ns/op\t39051910 B/op\t  134563 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42060403,
            "unit": "ns/op\t39053000 B/op\t  134580 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42305615,
            "unit": "ns/op\t39052695 B/op\t  134537 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55264187,
            "unit": "ns/op\t42619366 B/op\t  139461 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55304872,
            "unit": "ns/op\t42573635 B/op\t  139442 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55550314,
            "unit": "ns/op\t42542632 B/op\t  139413 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55149738,
            "unit": "ns/op\t42587481 B/op\t  139448 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55047218,
            "unit": "ns/op\t42577458 B/op\t  139463 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47494589,
            "unit": "ns/op\t41556244 B/op\t  156018 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46485926,
            "unit": "ns/op\t41545931 B/op\t  156014 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46638415,
            "unit": "ns/op\t41592674 B/op\t  156056 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46312021,
            "unit": "ns/op\t41562249 B/op\t  156027 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46779365,
            "unit": "ns/op\t41577215 B/op\t  156070 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46821369,
            "unit": "ns/op\t41600711 B/op\t  156108 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46400446,
            "unit": "ns/op\t41578983 B/op\t  156051 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47217549,
            "unit": "ns/op\t41560673 B/op\t  156039 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46907756,
            "unit": "ns/op\t41573325 B/op\t  156060 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46989648,
            "unit": "ns/op\t41548228 B/op\t  156039 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130803017,
            "unit": "ns/op\t60651992 B/op\t  715429 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129720908,
            "unit": "ns/op\t60823601 B/op\t  715539 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129073524,
            "unit": "ns/op\t60480651 B/op\t  715370 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129551076,
            "unit": "ns/op\t60567580 B/op\t  715433 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130787583,
            "unit": "ns/op\t60566805 B/op\t  715405 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "10d1c5433e10febc615e197a9a30263bd9530203",
          "message": "function: add support for time() (#141)\n\nAdd support for time() function. Add special handling for functions that\r\ntake no input and just return something for each step. It will be used\r\nin the future for other date functions.\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>\r\n\r\nSigned-off-by: Giedrius Statkevičius <giedrius.statkevicius@vinted.com>",
          "timestamp": "2023-01-05T15:08:48+02:00",
          "tree_id": "37f88314bf2679ecc85cb0179869f494d5292a1e",
          "url": "https://github.com/thanos-community/promql-engine/commit/10d1c5433e10febc615e197a9a30263bd9530203"
        },
        "date": 1672924298487,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25120411,
            "unit": "ns/op\t38751758 B/op\t  131449 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25412690,
            "unit": "ns/op\t38760546 B/op\t  131437 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25407087,
            "unit": "ns/op\t38736844 B/op\t  131435 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25986644,
            "unit": "ns/op\t38771284 B/op\t  131467 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25906417,
            "unit": "ns/op\t38806969 B/op\t  131484 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11984103,
            "unit": "ns/op\t12218817 B/op\t  126088 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12040552,
            "unit": "ns/op\t12234721 B/op\t  126079 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12019577,
            "unit": "ns/op\t12198481 B/op\t  126057 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12030758,
            "unit": "ns/op\t12210603 B/op\t  126071 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11999646,
            "unit": "ns/op\t12199003 B/op\t  126065 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21043077,
            "unit": "ns/op\t24263485 B/op\t  211842 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20798988,
            "unit": "ns/op\t24263532 B/op\t  211847 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20980956,
            "unit": "ns/op\t24271759 B/op\t  211845 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20942223,
            "unit": "ns/op\t24246782 B/op\t  211826 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 20996469,
            "unit": "ns/op\t24270336 B/op\t  211844 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16369593,
            "unit": "ns/op\t12005185 B/op\t  133422 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16322970,
            "unit": "ns/op\t12002132 B/op\t  133383 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16376552,
            "unit": "ns/op\t12060612 B/op\t  133493 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16333899,
            "unit": "ns/op\t11997980 B/op\t  133363 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16313132,
            "unit": "ns/op\t12031184 B/op\t  133505 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16350685,
            "unit": "ns/op\t11992639 B/op\t  131483 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16330072,
            "unit": "ns/op\t11956274 B/op\t  131432 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16321517,
            "unit": "ns/op\t12000212 B/op\t  131467 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16313290,
            "unit": "ns/op\t12012186 B/op\t  131473 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16373451,
            "unit": "ns/op\t11943809 B/op\t  131457 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33201886,
            "unit": "ns/op\t41596784 B/op\t  155337 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33052260,
            "unit": "ns/op\t41555076 B/op\t  155262 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33357180,
            "unit": "ns/op\t41607409 B/op\t  155308 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33133790,
            "unit": "ns/op\t41607138 B/op\t  155290 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33268318,
            "unit": "ns/op\t41570200 B/op\t  155273 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24085984,
            "unit": "ns/op\t14694278 B/op\t  149878 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24012609,
            "unit": "ns/op\t14732314 B/op\t  149886 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24111634,
            "unit": "ns/op\t14697200 B/op\t  149892 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23733380,
            "unit": "ns/op\t14693645 B/op\t  149871 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24114281,
            "unit": "ns/op\t14612937 B/op\t  149794 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32684546,
            "unit": "ns/op\t26050193 B/op\t  235190 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32941868,
            "unit": "ns/op\t26057163 B/op\t  235182 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32743477,
            "unit": "ns/op\t26044280 B/op\t  235176 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32884306,
            "unit": "ns/op\t26088322 B/op\t  235228 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32939269,
            "unit": "ns/op\t26105670 B/op\t  235231 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30181195,
            "unit": "ns/op\t44459900 B/op\t  619519 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30089528,
            "unit": "ns/op\t44430656 B/op\t  619526 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30116156,
            "unit": "ns/op\t44480899 B/op\t  619552 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30005860,
            "unit": "ns/op\t44472145 B/op\t  619522 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30129188,
            "unit": "ns/op\t44465308 B/op\t  619539 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18753794,
            "unit": "ns/op\t22943616 B/op\t  108842 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18939857,
            "unit": "ns/op\t22953493 B/op\t  108871 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18837857,
            "unit": "ns/op\t22930275 B/op\t  108826 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18851557,
            "unit": "ns/op\t22947534 B/op\t  108841 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18877274,
            "unit": "ns/op\t22934186 B/op\t  108835 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43456147,
            "unit": "ns/op\t57812642 B/op\t  202849 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44245622,
            "unit": "ns/op\t57836819 B/op\t  202862 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43951692,
            "unit": "ns/op\t57826258 B/op\t  202865 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44261670,
            "unit": "ns/op\t57836925 B/op\t  202882 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43978004,
            "unit": "ns/op\t57856580 B/op\t  202894 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37029122,
            "unit": "ns/op\t42272127 B/op\t  135756 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36560158,
            "unit": "ns/op\t42352025 B/op\t  135808 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36677828,
            "unit": "ns/op\t42192147 B/op\t  135727 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36560639,
            "unit": "ns/op\t42300739 B/op\t  135771 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36304445,
            "unit": "ns/op\t42233743 B/op\t  135755 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28646619,
            "unit": "ns/op\t39978944 B/op\t  143673 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28273816,
            "unit": "ns/op\t40039601 B/op\t  143714 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28187622,
            "unit": "ns/op\t40041976 B/op\t  143726 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28652515,
            "unit": "ns/op\t40022840 B/op\t  143734 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28222557,
            "unit": "ns/op\t40060723 B/op\t  143748 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37030155,
            "unit": "ns/op\t41994499 B/op\t  132806 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36616712,
            "unit": "ns/op\t41942120 B/op\t  132781 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36712995,
            "unit": "ns/op\t41960754 B/op\t  132783 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36976225,
            "unit": "ns/op\t41897181 B/op\t  132730 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36559242,
            "unit": "ns/op\t41850274 B/op\t  132725 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26488714,
            "unit": "ns/op\t36959699 B/op\t  102496 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26409037,
            "unit": "ns/op\t37003934 B/op\t  102521 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26052032,
            "unit": "ns/op\t36992732 B/op\t  102521 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26635507,
            "unit": "ns/op\t36968333 B/op\t  102500 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26548939,
            "unit": "ns/op\t36990253 B/op\t  102534 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21228558,
            "unit": "ns/op\t33144825 B/op\t   74815 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21389229,
            "unit": "ns/op\t33142786 B/op\t   74815 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21558086,
            "unit": "ns/op\t33144161 B/op\t   74815 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21741690,
            "unit": "ns/op\t33142485 B/op\t   74815 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21730262,
            "unit": "ns/op\t33142632 B/op\t   74815 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20632887,
            "unit": "ns/op\t32966152 B/op\t   68816 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20768415,
            "unit": "ns/op\t32964355 B/op\t   68816 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20368856,
            "unit": "ns/op\t32963756 B/op\t   68816 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21019812,
            "unit": "ns/op\t32963305 B/op\t   68816 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21218779,
            "unit": "ns/op\t32969629 B/op\t   68816 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44548511,
            "unit": "ns/op\t39150357 B/op\t  135242 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44354941,
            "unit": "ns/op\t39120178 B/op\t  135249 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44661190,
            "unit": "ns/op\t39149318 B/op\t  135256 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44343325,
            "unit": "ns/op\t39152839 B/op\t  135271 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44187766,
            "unit": "ns/op\t39153834 B/op\t  135296 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42007351,
            "unit": "ns/op\t39098777 B/op\t  134847 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41375286,
            "unit": "ns/op\t39085726 B/op\t  134811 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41740732,
            "unit": "ns/op\t39082896 B/op\t  134808 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41956732,
            "unit": "ns/op\t39095506 B/op\t  134817 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41286542,
            "unit": "ns/op\t39134481 B/op\t  134861 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54824712,
            "unit": "ns/op\t42547281 B/op\t  139662 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55036994,
            "unit": "ns/op\t42492871 B/op\t  139623 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54502092,
            "unit": "ns/op\t42483049 B/op\t  139609 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 53828092,
            "unit": "ns/op\t42439491 B/op\t  139589 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55055584,
            "unit": "ns/op\t42429604 B/op\t  139580 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46921559,
            "unit": "ns/op\t41575910 B/op\t  156274 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46760371,
            "unit": "ns/op\t41558824 B/op\t  156273 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46395171,
            "unit": "ns/op\t41582909 B/op\t  156253 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46647950,
            "unit": "ns/op\t41557937 B/op\t  156261 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46586051,
            "unit": "ns/op\t41548999 B/op\t  156213 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46981034,
            "unit": "ns/op\t41595294 B/op\t  156296 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46115944,
            "unit": "ns/op\t41633363 B/op\t  156385 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46624630,
            "unit": "ns/op\t41603810 B/op\t  156328 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46065594,
            "unit": "ns/op\t41563664 B/op\t  156254 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46708039,
            "unit": "ns/op\t41585947 B/op\t  156292 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130044484,
            "unit": "ns/op\t60750035 B/op\t  715305 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130789363,
            "unit": "ns/op\t60473373 B/op\t  715142 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129885952,
            "unit": "ns/op\t60758714 B/op\t  715201 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129027176,
            "unit": "ns/op\t60453177 B/op\t  715090 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130364550,
            "unit": "ns/op\t60592626 B/op\t  715209 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "giedrius.statkevicius@vinted.com",
            "name": "Giedrius Statkevičius",
            "username": "GiedriusS"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a831c7b187d4ae915e625aba3be45607735003e1",
          "message": "function: add support for more functions (#142)\n\nAdd various math functions and tests for them. Fix fuzzing test to\r\nignore functions that do not consume a matrix.",
          "timestamp": "2023-01-06T10:45:17+02:00",
          "tree_id": "1f2141ea1e457dd3b55820e62cdf49ed7de1df94",
          "url": "https://github.com/thanos-community/promql-engine/commit/a831c7b187d4ae915e625aba3be45607735003e1"
        },
        "date": 1672994886021,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25880142,
            "unit": "ns/op\t38762534 B/op\t  131229 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25909054,
            "unit": "ns/op\t38739917 B/op\t  131194 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26641960,
            "unit": "ns/op\t38729131 B/op\t  131188 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26561664,
            "unit": "ns/op\t38742534 B/op\t  131204 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26758785,
            "unit": "ns/op\t38742920 B/op\t  131201 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11774446,
            "unit": "ns/op\t12189524 B/op\t  125822 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11814918,
            "unit": "ns/op\t12190235 B/op\t  125815 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11788911,
            "unit": "ns/op\t12227413 B/op\t  125840 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11757523,
            "unit": "ns/op\t12179123 B/op\t  125819 allocs/op",
            "extra": "99 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11786719,
            "unit": "ns/op\t12195354 B/op\t  125836 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21269823,
            "unit": "ns/op\t24229439 B/op\t  211563 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21187473,
            "unit": "ns/op\t24235384 B/op\t  211569 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21101931,
            "unit": "ns/op\t24242515 B/op\t  211572 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21115723,
            "unit": "ns/op\t24257673 B/op\t  211590 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21270807,
            "unit": "ns/op\t24254398 B/op\t  211584 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16515327,
            "unit": "ns/op\t12177597 B/op\t  133470 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16407185,
            "unit": "ns/op\t12105542 B/op\t  133368 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16138445,
            "unit": "ns/op\t12063074 B/op\t  133281 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16277116,
            "unit": "ns/op\t12137466 B/op\t  133443 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16288083,
            "unit": "ns/op\t12146844 B/op\t  133449 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16343835,
            "unit": "ns/op\t12214605 B/op\t  134303 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16332228,
            "unit": "ns/op\t12172188 B/op\t  134294 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16256986,
            "unit": "ns/op\t12136036 B/op\t  134240 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16301330,
            "unit": "ns/op\t12139195 B/op\t  134279 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16341757,
            "unit": "ns/op\t12177584 B/op\t  134298 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33458764,
            "unit": "ns/op\t41560309 B/op\t  155005 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33494579,
            "unit": "ns/op\t41573459 B/op\t  155073 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33256502,
            "unit": "ns/op\t41540560 B/op\t  155041 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33748985,
            "unit": "ns/op\t41548930 B/op\t  155026 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33378201,
            "unit": "ns/op\t41570582 B/op\t  155063 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 23897526,
            "unit": "ns/op\t14700108 B/op\t  149628 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24062528,
            "unit": "ns/op\t14683759 B/op\t  149632 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24039869,
            "unit": "ns/op\t14650609 B/op\t  149578 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24045907,
            "unit": "ns/op\t14670034 B/op\t  149586 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24001608,
            "unit": "ns/op\t14674179 B/op\t  149608 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32996405,
            "unit": "ns/op\t26098819 B/op\t  234962 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32901198,
            "unit": "ns/op\t26172345 B/op\t  235030 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32654388,
            "unit": "ns/op\t26135521 B/op\t  234978 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32701378,
            "unit": "ns/op\t26087350 B/op\t  234965 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 32817464,
            "unit": "ns/op\t26092089 B/op\t  234944 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30363946,
            "unit": "ns/op\t44422855 B/op\t  619014 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30228539,
            "unit": "ns/op\t44410242 B/op\t  619035 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30220572,
            "unit": "ns/op\t44481305 B/op\t  619026 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30330401,
            "unit": "ns/op\t44400469 B/op\t  619045 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30037095,
            "unit": "ns/op\t44454479 B/op\t  619030 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18929856,
            "unit": "ns/op\t22919657 B/op\t  108439 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18949584,
            "unit": "ns/op\t22937014 B/op\t  108484 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18939484,
            "unit": "ns/op\t22929787 B/op\t  108455 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18938366,
            "unit": "ns/op\t22923998 B/op\t  108459 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18995394,
            "unit": "ns/op\t22929649 B/op\t  108473 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44257382,
            "unit": "ns/op\t57840094 B/op\t  202304 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44097673,
            "unit": "ns/op\t57854643 B/op\t  202313 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44427513,
            "unit": "ns/op\t57842971 B/op\t  202316 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44275844,
            "unit": "ns/op\t57787459 B/op\t  202265 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44302284,
            "unit": "ns/op\t57854839 B/op\t  202336 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37216978,
            "unit": "ns/op\t42393053 B/op\t  135582 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37430769,
            "unit": "ns/op\t42248510 B/op\t  135502 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36978462,
            "unit": "ns/op\t42379855 B/op\t  135579 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37225753,
            "unit": "ns/op\t42263009 B/op\t  135510 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36737912,
            "unit": "ns/op\t42372542 B/op\t  135564 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28359029,
            "unit": "ns/op\t40016637 B/op\t  143490 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28697466,
            "unit": "ns/op\t40001208 B/op\t  143470 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28496459,
            "unit": "ns/op\t40001737 B/op\t  143466 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28552359,
            "unit": "ns/op\t40016543 B/op\t  143472 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28724104,
            "unit": "ns/op\t40026109 B/op\t  143496 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37330755,
            "unit": "ns/op\t41904163 B/op\t  132510 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36937971,
            "unit": "ns/op\t41921509 B/op\t  132506 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37395738,
            "unit": "ns/op\t41860527 B/op\t  132482 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36790317,
            "unit": "ns/op\t41904286 B/op\t  132518 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36852078,
            "unit": "ns/op\t41933500 B/op\t  132530 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26789985,
            "unit": "ns/op\t36980420 B/op\t  102350 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26822156,
            "unit": "ns/op\t36949102 B/op\t  102333 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26535047,
            "unit": "ns/op\t37006073 B/op\t  102374 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26897968,
            "unit": "ns/op\t36998549 B/op\t  102359 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26704733,
            "unit": "ns/op\t36976982 B/op\t  102358 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22098682,
            "unit": "ns/op\t33138225 B/op\t   74730 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22379172,
            "unit": "ns/op\t33145429 B/op\t   74732 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21508071,
            "unit": "ns/op\t33141418 B/op\t   74731 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21427952,
            "unit": "ns/op\t33139246 B/op\t   74731 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22523347,
            "unit": "ns/op\t33141977 B/op\t   74732 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20846038,
            "unit": "ns/op\t32964649 B/op\t   68733 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21418249,
            "unit": "ns/op\t32959255 B/op\t   68732 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20923811,
            "unit": "ns/op\t32960246 B/op\t   68733 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21027980,
            "unit": "ns/op\t32959196 B/op\t   68733 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20715389,
            "unit": "ns/op\t32960986 B/op\t   68733 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44633167,
            "unit": "ns/op\t39088916 B/op\t  134973 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45199400,
            "unit": "ns/op\t39107341 B/op\t  134972 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44441180,
            "unit": "ns/op\t39111853 B/op\t  135010 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44959929,
            "unit": "ns/op\t39124806 B/op\t  134999 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44248945,
            "unit": "ns/op\t39123121 B/op\t  134999 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42287697,
            "unit": "ns/op\t39056382 B/op\t  134533 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41306772,
            "unit": "ns/op\t39071992 B/op\t  134595 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 40673732,
            "unit": "ns/op\t39099108 B/op\t  134602 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42526195,
            "unit": "ns/op\t39055936 B/op\t  134547 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41297928,
            "unit": "ns/op\t39085102 B/op\t  134587 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55814847,
            "unit": "ns/op\t42537120 B/op\t  139418 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55669348,
            "unit": "ns/op\t42568129 B/op\t  139427 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54597732,
            "unit": "ns/op\t42504827 B/op\t  139398 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 54667196,
            "unit": "ns/op\t42492904 B/op\t  139392 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55304325,
            "unit": "ns/op\t42449102 B/op\t  139355 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46698585,
            "unit": "ns/op\t41588016 B/op\t  156062 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47645576,
            "unit": "ns/op\t41532712 B/op\t  156011 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46962207,
            "unit": "ns/op\t41574152 B/op\t  156050 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46925584,
            "unit": "ns/op\t41587452 B/op\t  156076 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 46829189,
            "unit": "ns/op\t41561328 B/op\t  156034 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47048083,
            "unit": "ns/op\t41551559 B/op\t  156008 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47464133,
            "unit": "ns/op\t41570150 B/op\t  156050 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47145731,
            "unit": "ns/op\t41538884 B/op\t  156017 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46676584,
            "unit": "ns/op\t41565904 B/op\t  156082 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47313427,
            "unit": "ns/op\t41558066 B/op\t  156036 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129452075,
            "unit": "ns/op\t60695854 B/op\t  715701 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 127975788,
            "unit": "ns/op\t60517763 B/op\t  715595 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129140420,
            "unit": "ns/op\t60739087 B/op\t  715667 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 127255503,
            "unit": "ns/op\t60707845 B/op\t  715648 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130803518,
            "unit": "ns/op\t60698010 B/op\t  715659 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "zhaoziqi9146@gmail.com",
            "name": "Ziqi Zhao",
            "username": "fatsheep9146"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3773ec307b6a1a3aa9ccada12608df50d625de81",
          "message": "function: add support for timestamp() (#145)\n\nSigned-off-by: Ziqi Zhao <zhaoziqi9146@gmail.com>\r\n\r\nSigned-off-by: Ziqi Zhao <zhaoziqi9146@gmail.com>",
          "timestamp": "2023-01-09T07:34:14+01:00",
          "tree_id": "5b783495bb6daf968b2eddeb4b1d5879b8150ee2",
          "url": "https://github.com/thanos-community/promql-engine/commit/3773ec307b6a1a3aa9ccada12608df50d625de81"
        },
        "date": 1673246225129,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25712781,
            "unit": "ns/op\t38761132 B/op\t  131493 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25451967,
            "unit": "ns/op\t38794156 B/op\t  131510 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 25661522,
            "unit": "ns/op\t38746572 B/op\t  131466 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 26255385,
            "unit": "ns/op\t38759595 B/op\t  131478 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 27060097,
            "unit": "ns/op\t38741357 B/op\t  131474 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12105203,
            "unit": "ns/op\t12204757 B/op\t  126125 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12082688,
            "unit": "ns/op\t12246926 B/op\t  126134 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12060482,
            "unit": "ns/op\t12219559 B/op\t  126118 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12103327,
            "unit": "ns/op\t12187080 B/op\t  126094 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12056221,
            "unit": "ns/op\t12231005 B/op\t  126123 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21185050,
            "unit": "ns/op\t24261987 B/op\t  211864 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21159719,
            "unit": "ns/op\t24283607 B/op\t  211866 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21126028,
            "unit": "ns/op\t24285199 B/op\t  211878 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21335767,
            "unit": "ns/op\t24262480 B/op\t  211867 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 21286841,
            "unit": "ns/op\t24270709 B/op\t  211866 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16538499,
            "unit": "ns/op\t12023538 B/op\t  130593 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16495642,
            "unit": "ns/op\t12010768 B/op\t  130557 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16548704,
            "unit": "ns/op\t12020218 B/op\t  130622 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16548388,
            "unit": "ns/op\t12000311 B/op\t  130589 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 16429619,
            "unit": "ns/op\t11977802 B/op\t  130553 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16573132,
            "unit": "ns/op\t12019402 B/op\t  131771 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16513308,
            "unit": "ns/op\t12023401 B/op\t  131707 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16430711,
            "unit": "ns/op\t12012044 B/op\t  131690 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16452563,
            "unit": "ns/op\t11971080 B/op\t  131578 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 16446466,
            "unit": "ns/op\t12038063 B/op\t  131738 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33849030,
            "unit": "ns/op\t41545047 B/op\t  155302 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33447715,
            "unit": "ns/op\t41580828 B/op\t  155316 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33790838,
            "unit": "ns/op\t41613569 B/op\t  155343 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33707626,
            "unit": "ns/op\t41574868 B/op\t  155310 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 33440894,
            "unit": "ns/op\t41590724 B/op\t  155314 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24018909,
            "unit": "ns/op\t14777037 B/op\t  149966 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24083713,
            "unit": "ns/op\t14691122 B/op\t  149902 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24259747,
            "unit": "ns/op\t14631952 B/op\t  149848 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24153028,
            "unit": "ns/op\t14707106 B/op\t  149899 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24191264,
            "unit": "ns/op\t14710387 B/op\t  149898 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33116253,
            "unit": "ns/op\t26123144 B/op\t  235257 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33100836,
            "unit": "ns/op\t26107041 B/op\t  235234 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33288592,
            "unit": "ns/op\t26087994 B/op\t  235240 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33035359,
            "unit": "ns/op\t26082274 B/op\t  235229 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 33065954,
            "unit": "ns/op\t26221376 B/op\t  235317 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30415879,
            "unit": "ns/op\t44508248 B/op\t  619603 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30507942,
            "unit": "ns/op\t44444567 B/op\t  619598 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30538865,
            "unit": "ns/op\t44454617 B/op\t  619592 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30680603,
            "unit": "ns/op\t44454462 B/op\t  619575 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 30346051,
            "unit": "ns/op\t44434701 B/op\t  619565 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18860222,
            "unit": "ns/op\t22941385 B/op\t  108546 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19084832,
            "unit": "ns/op\t22934748 B/op\t  108559 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19044630,
            "unit": "ns/op\t22927336 B/op\t  108525 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19056310,
            "unit": "ns/op\t22948251 B/op\t  108578 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18913440,
            "unit": "ns/op\t22932356 B/op\t  108538 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44550651,
            "unit": "ns/op\t57827218 B/op\t  202449 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44405888,
            "unit": "ns/op\t57870091 B/op\t  202478 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44108732,
            "unit": "ns/op\t57805578 B/op\t  202415 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44372400,
            "unit": "ns/op\t57858349 B/op\t  202462 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 44158646,
            "unit": "ns/op\t57763157 B/op\t  202365 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 37226177,
            "unit": "ns/op\t42222159 B/op\t  135759 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36841267,
            "unit": "ns/op\t42293800 B/op\t  135809 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36375189,
            "unit": "ns/op\t42278023 B/op\t  135812 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36882473,
            "unit": "ns/op\t42332728 B/op\t  135829 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 36842387,
            "unit": "ns/op\t42407553 B/op\t  135865 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28687473,
            "unit": "ns/op\t40006376 B/op\t  143724 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28897499,
            "unit": "ns/op\t40000928 B/op\t  143719 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28984581,
            "unit": "ns/op\t40035692 B/op\t  143750 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28659619,
            "unit": "ns/op\t40015174 B/op\t  143747 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 28386069,
            "unit": "ns/op\t40016465 B/op\t  143728 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36653926,
            "unit": "ns/op\t41988894 B/op\t  132849 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37251779,
            "unit": "ns/op\t41891946 B/op\t  132776 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 36887965,
            "unit": "ns/op\t41954718 B/op\t  132804 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37043426,
            "unit": "ns/op\t41926536 B/op\t  132779 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 37454448,
            "unit": "ns/op\t41923577 B/op\t  132792 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26828822,
            "unit": "ns/op\t36979703 B/op\t  102545 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26444365,
            "unit": "ns/op\t37031439 B/op\t  102582 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26449672,
            "unit": "ns/op\t37017974 B/op\t  102566 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26819417,
            "unit": "ns/op\t36976971 B/op\t  102531 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 26793894,
            "unit": "ns/op\t37014940 B/op\t  102563 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21613081,
            "unit": "ns/op\t33146618 B/op\t   74827 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22346280,
            "unit": "ns/op\t33142943 B/op\t   74825 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22257394,
            "unit": "ns/op\t33143519 B/op\t   74827 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 22454040,
            "unit": "ns/op\t33149482 B/op\t   74828 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 21901679,
            "unit": "ns/op\t33143295 B/op\t   74827 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20466693,
            "unit": "ns/op\t32964632 B/op\t   68828 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 20838604,
            "unit": "ns/op\t32966554 B/op\t   68829 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21204084,
            "unit": "ns/op\t32968178 B/op\t   68828 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21034301,
            "unit": "ns/op\t32970093 B/op\t   68829 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 21257133,
            "unit": "ns/op\t32973512 B/op\t   68830 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44412238,
            "unit": "ns/op\t39133724 B/op\t  135282 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44655647,
            "unit": "ns/op\t39134298 B/op\t  135303 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45115009,
            "unit": "ns/op\t39146131 B/op\t  135275 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 44499129,
            "unit": "ns/op\t39126604 B/op\t  135241 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 45067797,
            "unit": "ns/op\t39129781 B/op\t  135292 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42079136,
            "unit": "ns/op\t39074020 B/op\t  134828 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 42139080,
            "unit": "ns/op\t39063288 B/op\t  134842 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41299025,
            "unit": "ns/op\t39121994 B/op\t  134886 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41924317,
            "unit": "ns/op\t39076114 B/op\t  134831 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 41619906,
            "unit": "ns/op\t39113755 B/op\t  134884 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55932784,
            "unit": "ns/op\t42470218 B/op\t  139633 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55478037,
            "unit": "ns/op\t42393554 B/op\t  139559 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55014952,
            "unit": "ns/op\t42406995 B/op\t  139603 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55151582,
            "unit": "ns/op\t42614922 B/op\t  139726 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 55278399,
            "unit": "ns/op\t42384219 B/op\t  139574 allocs/op",
            "extra": "21 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47228129,
            "unit": "ns/op\t41579792 B/op\t  156313 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47623977,
            "unit": "ns/op\t41579124 B/op\t  156319 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47444712,
            "unit": "ns/op\t41562002 B/op\t  156281 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47308523,
            "unit": "ns/op\t41578873 B/op\t  156320 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 47107947,
            "unit": "ns/op\t41623793 B/op\t  156378 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47170400,
            "unit": "ns/op\t41614054 B/op\t  156360 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47519728,
            "unit": "ns/op\t41576160 B/op\t  156312 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 46958067,
            "unit": "ns/op\t41597417 B/op\t  156342 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47233726,
            "unit": "ns/op\t41537571 B/op\t  156269 allocs/op",
            "extra": "25 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 47014873,
            "unit": "ns/op\t41595013 B/op\t  156312 allocs/op",
            "extra": "26 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129668387,
            "unit": "ns/op\t60494487 B/op\t  715523 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130344876,
            "unit": "ns/op\t60787804 B/op\t  715636 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 130003316,
            "unit": "ns/op\t60471286 B/op\t  715501 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129499759,
            "unit": "ns/op\t60624937 B/op\t  715498 allocs/op",
            "extra": "8 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 129350366,
            "unit": "ns/op\t60464654 B/op\t  715505 allocs/op",
            "extra": "8 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c293f65f538929eab77c6159936be35064536890",
          "message": "Implement distributed execution (#139)\n\n* Implement distributed query execution\r\n\r\n* Add test case for unsupported ops\r\n\r\n* Optimize point selection\r\n\r\n* Add concurrency to remote exec operator\r\n\r\n* Simplify remote storage\r\n\r\n* Improve tests\r\n\r\n* Allow engines to be discovered\r\n\r\n* Fix tests\r\n\r\n* Fix lint\r\n\r\n* Fix tests",
          "timestamp": "2023-01-11T13:15:31+01:00",
          "tree_id": "9f08ac7b71b7c5fee2e8fb3eb19fb0c44c2238c2",
          "url": "https://github.com/thanos-community/promql-engine/commit/c293f65f538929eab77c6159936be35064536890"
        },
        "date": 1673439503906,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22364719,
            "unit": "ns/op\t38809966 B/op\t  131463 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22251132,
            "unit": "ns/op\t38849338 B/op\t  131460 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22437518,
            "unit": "ns/op\t38851378 B/op\t  131469 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22889016,
            "unit": "ns/op\t38779818 B/op\t  131413 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22705485,
            "unit": "ns/op\t38823206 B/op\t  131461 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12019944,
            "unit": "ns/op\t10950669 B/op\t  125978 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12068027,
            "unit": "ns/op\t10940582 B/op\t  125982 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12073273,
            "unit": "ns/op\t10949812 B/op\t  125985 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12085170,
            "unit": "ns/op\t10979869 B/op\t  126022 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12075376,
            "unit": "ns/op\t10956531 B/op\t  126001 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17099794,
            "unit": "ns/op\t24350378 B/op\t  211813 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17140656,
            "unit": "ns/op\t24345754 B/op\t  211798 allocs/op",
            "extra": "74 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17178913,
            "unit": "ns/op\t24385089 B/op\t  211839 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17257763,
            "unit": "ns/op\t24370859 B/op\t  211821 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17135879,
            "unit": "ns/op\t24373666 B/op\t  211821 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12424152,
            "unit": "ns/op\t12335923 B/op\t  134038 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12400843,
            "unit": "ns/op\t12298462 B/op\t  134008 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12462464,
            "unit": "ns/op\t12324954 B/op\t  134024 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12454762,
            "unit": "ns/op\t12345965 B/op\t  134087 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12399719,
            "unit": "ns/op\t12319828 B/op\t  133996 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12401939,
            "unit": "ns/op\t12277449 B/op\t  132613 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12405019,
            "unit": "ns/op\t12299314 B/op\t  132658 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12408722,
            "unit": "ns/op\t12293978 B/op\t  132600 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12391896,
            "unit": "ns/op\t12340796 B/op\t  132679 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12448785,
            "unit": "ns/op\t12290085 B/op\t  132656 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28542577,
            "unit": "ns/op\t41762780 B/op\t  155338 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28509468,
            "unit": "ns/op\t41731242 B/op\t  155321 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28644662,
            "unit": "ns/op\t41754550 B/op\t  155334 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28582312,
            "unit": "ns/op\t41731642 B/op\t  155294 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28601228,
            "unit": "ns/op\t41742200 B/op\t  155326 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24571996,
            "unit": "ns/op\t13555340 B/op\t  149700 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24683588,
            "unit": "ns/op\t13563650 B/op\t  149716 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24535003,
            "unit": "ns/op\t13652374 B/op\t  149802 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24441572,
            "unit": "ns/op\t13591153 B/op\t  149754 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24616584,
            "unit": "ns/op\t13584634 B/op\t  149735 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27740222,
            "unit": "ns/op\t26270493 B/op\t  235249 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27794283,
            "unit": "ns/op\t26221507 B/op\t  235206 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27849090,
            "unit": "ns/op\t26206978 B/op\t  235204 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27777200,
            "unit": "ns/op\t26289416 B/op\t  235277 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27840476,
            "unit": "ns/op\t26233606 B/op\t  235231 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26626466,
            "unit": "ns/op\t42548699 B/op\t  619291 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26306299,
            "unit": "ns/op\t42582779 B/op\t  619312 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26329937,
            "unit": "ns/op\t42505607 B/op\t  619291 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26379351,
            "unit": "ns/op\t42522942 B/op\t  619278 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26434602,
            "unit": "ns/op\t42563566 B/op\t  619306 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17888080,
            "unit": "ns/op\t22982451 B/op\t  108552 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18000845,
            "unit": "ns/op\t22990930 B/op\t  108565 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17916192,
            "unit": "ns/op\t22989698 B/op\t  108558 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18100992,
            "unit": "ns/op\t22998244 B/op\t  108574 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17835505,
            "unit": "ns/op\t22986357 B/op\t  108561 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40005682,
            "unit": "ns/op\t57892580 B/op\t  202493 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39884195,
            "unit": "ns/op\t57929724 B/op\t  202537 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39734076,
            "unit": "ns/op\t57923301 B/op\t  202536 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39538213,
            "unit": "ns/op\t57914622 B/op\t  202496 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39455953,
            "unit": "ns/op\t57928744 B/op\t  202560 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33054742,
            "unit": "ns/op\t42414068 B/op\t  135776 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32976313,
            "unit": "ns/op\t42505066 B/op\t  135822 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32352195,
            "unit": "ns/op\t42387716 B/op\t  135768 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32268335,
            "unit": "ns/op\t42489475 B/op\t  135840 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32629936,
            "unit": "ns/op\t42399588 B/op\t  135773 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24202630,
            "unit": "ns/op\t40111924 B/op\t  143746 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24548089,
            "unit": "ns/op\t40125200 B/op\t  143749 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24095512,
            "unit": "ns/op\t40125581 B/op\t  143744 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24470421,
            "unit": "ns/op\t40091111 B/op\t  143731 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24428509,
            "unit": "ns/op\t40094445 B/op\t  143722 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32730088,
            "unit": "ns/op\t41975181 B/op\t  132742 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32711946,
            "unit": "ns/op\t42077946 B/op\t  132788 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32946285,
            "unit": "ns/op\t41990273 B/op\t  132741 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32819808,
            "unit": "ns/op\t42089019 B/op\t  132798 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32819599,
            "unit": "ns/op\t42043068 B/op\t  132783 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 21968501,
            "unit": "ns/op\t37049836 B/op\t  102520 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22301117,
            "unit": "ns/op\t37042886 B/op\t  102516 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22083018,
            "unit": "ns/op\t37066623 B/op\t  102535 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22097620,
            "unit": "ns/op\t37035618 B/op\t  102515 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22064328,
            "unit": "ns/op\t37074198 B/op\t  102542 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 18009592,
            "unit": "ns/op\t33220028 B/op\t   74834 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17307063,
            "unit": "ns/op\t33218866 B/op\t   74835 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17301782,
            "unit": "ns/op\t33222325 B/op\t   74834 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17069024,
            "unit": "ns/op\t33221600 B/op\t   74835 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17544188,
            "unit": "ns/op\t33222386 B/op\t   74834 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17406342,
            "unit": "ns/op\t33045736 B/op\t   68837 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16636586,
            "unit": "ns/op\t33043082 B/op\t   68837 allocs/op",
            "extra": "86 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16956595,
            "unit": "ns/op\t33041240 B/op\t   68836 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17163687,
            "unit": "ns/op\t33044448 B/op\t   68836 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16881158,
            "unit": "ns/op\t33046072 B/op\t   68838 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40431643,
            "unit": "ns/op\t39143968 B/op\t  135188 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40504197,
            "unit": "ns/op\t39189092 B/op\t  135227 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40122169,
            "unit": "ns/op\t39170856 B/op\t  135221 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40412181,
            "unit": "ns/op\t39178672 B/op\t  135239 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40449534,
            "unit": "ns/op\t39167634 B/op\t  135222 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38043538,
            "unit": "ns/op\t39150323 B/op\t  134816 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37402633,
            "unit": "ns/op\t39127093 B/op\t  134800 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37525860,
            "unit": "ns/op\t39126375 B/op\t  134811 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37589399,
            "unit": "ns/op\t39115846 B/op\t  134779 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37608829,
            "unit": "ns/op\t39141555 B/op\t  134831 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51137025,
            "unit": "ns/op\t42596029 B/op\t  139625 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50756144,
            "unit": "ns/op\t42599602 B/op\t  139614 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51672354,
            "unit": "ns/op\t42556090 B/op\t  139602 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50219594,
            "unit": "ns/op\t42622995 B/op\t  139644 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50386359,
            "unit": "ns/op\t42636670 B/op\t  139641 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41840690,
            "unit": "ns/op\t41705587 B/op\t  156316 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41769127,
            "unit": "ns/op\t41701843 B/op\t  156282 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41895263,
            "unit": "ns/op\t41723995 B/op\t  156322 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42236940,
            "unit": "ns/op\t41693886 B/op\t  156297 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41771926,
            "unit": "ns/op\t41713376 B/op\t  156298 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41765056,
            "unit": "ns/op\t41677988 B/op\t  156269 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41731717,
            "unit": "ns/op\t41715399 B/op\t  156294 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42156069,
            "unit": "ns/op\t41686446 B/op\t  156281 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41794324,
            "unit": "ns/op\t41706017 B/op\t  156297 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41625530,
            "unit": "ns/op\t41749357 B/op\t  156318 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113867057,
            "unit": "ns/op\t61215677 B/op\t  715640 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114165002,
            "unit": "ns/op\t61224121 B/op\t  715675 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113649249,
            "unit": "ns/op\t61334210 B/op\t  715717 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114243685,
            "unit": "ns/op\t61058753 B/op\t  715602 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113880237,
            "unit": "ns/op\t61227971 B/op\t  715609 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benye@amazon.com",
            "name": "Ben Ye",
            "username": "yeya24"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3e3e0b3c510bae274b2af7bcf74eea40282c0425",
          "message": "Avoid panic in query stats (#153)\n\n* implement query stats to avoid panic\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>\r\n\r\n* fix lint\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>\r\n\r\nSigned-off-by: Ben Ye <benye@amazon.com>",
          "timestamp": "2023-01-17T09:20:09-08:00",
          "tree_id": "2613e7039f377e54324b31193feffcca0d288725",
          "url": "https://github.com/thanos-community/promql-engine/commit/3e3e0b3c510bae274b2af7bcf74eea40282c0425"
        },
        "date": 1673976181039,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 21690536,
            "unit": "ns/op\t38816696 B/op\t  131508 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 21756307,
            "unit": "ns/op\t38803585 B/op\t  131460 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 21831775,
            "unit": "ns/op\t38781017 B/op\t  131466 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 21611462,
            "unit": "ns/op\t38819457 B/op\t  131480 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22451967,
            "unit": "ns/op\t38789210 B/op\t  131462 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11868490,
            "unit": "ns/op\t10927364 B/op\t  125987 allocs/op",
            "extra": "100 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11975914,
            "unit": "ns/op\t10930549 B/op\t  126008 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12026113,
            "unit": "ns/op\t10928237 B/op\t  125995 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12013820,
            "unit": "ns/op\t10943810 B/op\t  126001 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 11991361,
            "unit": "ns/op\t10975249 B/op\t  126046 allocs/op",
            "extra": "99 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 16843148,
            "unit": "ns/op\t24354634 B/op\t  211849 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 16959094,
            "unit": "ns/op\t24381618 B/op\t  211856 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17012200,
            "unit": "ns/op\t24362417 B/op\t  211855 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17107566,
            "unit": "ns/op\t24362739 B/op\t  211860 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 16913260,
            "unit": "ns/op\t24371506 B/op\t  211856 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12341447,
            "unit": "ns/op\t12329904 B/op\t  134154 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12332875,
            "unit": "ns/op\t12301227 B/op\t  134119 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12386966,
            "unit": "ns/op\t12349745 B/op\t  134193 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12440757,
            "unit": "ns/op\t12374623 B/op\t  134221 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12306770,
            "unit": "ns/op\t12284279 B/op\t  134079 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12427736,
            "unit": "ns/op\t12411908 B/op\t  135529 allocs/op",
            "extra": "98 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12421760,
            "unit": "ns/op\t12400275 B/op\t  135473 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12369897,
            "unit": "ns/op\t12333633 B/op\t  135426 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12343692,
            "unit": "ns/op\t12373717 B/op\t  135489 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12373127,
            "unit": "ns/op\t12357475 B/op\t  135368 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28587045,
            "unit": "ns/op\t41752591 B/op\t  155315 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28672234,
            "unit": "ns/op\t41750874 B/op\t  155327 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28582775,
            "unit": "ns/op\t41763962 B/op\t  155347 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28321001,
            "unit": "ns/op\t41760752 B/op\t  155341 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28425207,
            "unit": "ns/op\t41753372 B/op\t  155345 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24361593,
            "unit": "ns/op\t13575164 B/op\t  149753 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24537105,
            "unit": "ns/op\t13633111 B/op\t  149804 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24536552,
            "unit": "ns/op\t13613461 B/op\t  149782 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24285169,
            "unit": "ns/op\t13638841 B/op\t  149837 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24416776,
            "unit": "ns/op\t13611229 B/op\t  149798 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27642534,
            "unit": "ns/op\t26191215 B/op\t  235222 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27685763,
            "unit": "ns/op\t26237237 B/op\t  235259 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27737257,
            "unit": "ns/op\t26179661 B/op\t  235228 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27600605,
            "unit": "ns/op\t26213483 B/op\t  235252 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27730534,
            "unit": "ns/op\t26293766 B/op\t  235299 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26288660,
            "unit": "ns/op\t42569725 B/op\t  619383 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26046493,
            "unit": "ns/op\t42554375 B/op\t  619367 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26102823,
            "unit": "ns/op\t42584874 B/op\t  619385 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26197599,
            "unit": "ns/op\t42539842 B/op\t  619353 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26131621,
            "unit": "ns/op\t42543013 B/op\t  619387 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17794155,
            "unit": "ns/op\t22977657 B/op\t  108477 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17858626,
            "unit": "ns/op\t22976763 B/op\t  108485 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17832026,
            "unit": "ns/op\t22974112 B/op\t  108471 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17963844,
            "unit": "ns/op\t22976640 B/op\t  108500 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17838885,
            "unit": "ns/op\t22979024 B/op\t  108498 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39807793,
            "unit": "ns/op\t57905759 B/op\t  202576 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39213322,
            "unit": "ns/op\t57895137 B/op\t  202585 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39158156,
            "unit": "ns/op\t57908158 B/op\t  202572 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39365698,
            "unit": "ns/op\t57909921 B/op\t  202578 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39678045,
            "unit": "ns/op\t57925557 B/op\t  202606 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32604767,
            "unit": "ns/op\t42409087 B/op\t  135811 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32491939,
            "unit": "ns/op\t42405171 B/op\t  135819 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32705657,
            "unit": "ns/op\t42521861 B/op\t  135868 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32547147,
            "unit": "ns/op\t42450477 B/op\t  135840 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32608541,
            "unit": "ns/op\t42497776 B/op\t  135857 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24365959,
            "unit": "ns/op\t40086647 B/op\t  143757 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24201179,
            "unit": "ns/op\t40101530 B/op\t  143760 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24369098,
            "unit": "ns/op\t40104280 B/op\t  143772 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24274452,
            "unit": "ns/op\t40122364 B/op\t  143786 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24388508,
            "unit": "ns/op\t40073882 B/op\t  143752 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32661201,
            "unit": "ns/op\t41953420 B/op\t  132762 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32451520,
            "unit": "ns/op\t41950685 B/op\t  132761 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32641927,
            "unit": "ns/op\t42105401 B/op\t  132847 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32569220,
            "unit": "ns/op\t42038930 B/op\t  132822 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32700006,
            "unit": "ns/op\t42094620 B/op\t  132837 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22311711,
            "unit": "ns/op\t37061779 B/op\t  102555 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22023265,
            "unit": "ns/op\t37050707 B/op\t  102547 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22086141,
            "unit": "ns/op\t37053986 B/op\t  102552 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22140553,
            "unit": "ns/op\t37018265 B/op\t  102526 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 21886006,
            "unit": "ns/op\t37086389 B/op\t  102579 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16830660,
            "unit": "ns/op\t33222683 B/op\t   74847 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16957144,
            "unit": "ns/op\t33221194 B/op\t   74846 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16787317,
            "unit": "ns/op\t33222466 B/op\t   74847 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17097614,
            "unit": "ns/op\t33223863 B/op\t   74847 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17677802,
            "unit": "ns/op\t33224661 B/op\t   74847 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16931965,
            "unit": "ns/op\t33046191 B/op\t   68850 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17399302,
            "unit": "ns/op\t33046191 B/op\t   68849 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16656179,
            "unit": "ns/op\t33045298 B/op\t   68849 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17240584,
            "unit": "ns/op\t33040428 B/op\t   68849 allocs/op",
            "extra": "76 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17035641,
            "unit": "ns/op\t33044028 B/op\t   68848 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40155833,
            "unit": "ns/op\t39176268 B/op\t  135250 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40364001,
            "unit": "ns/op\t39217832 B/op\t  135301 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40321044,
            "unit": "ns/op\t39211409 B/op\t  135294 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40193895,
            "unit": "ns/op\t39183948 B/op\t  135270 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40709009,
            "unit": "ns/op\t39164045 B/op\t  135242 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38095500,
            "unit": "ns/op\t39113197 B/op\t  134801 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37298935,
            "unit": "ns/op\t39185443 B/op\t  134871 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37302708,
            "unit": "ns/op\t39144781 B/op\t  134841 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37936099,
            "unit": "ns/op\t39153830 B/op\t  134851 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37848342,
            "unit": "ns/op\t39158976 B/op\t  134862 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51195170,
            "unit": "ns/op\t42578874 B/op\t  139648 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50655252,
            "unit": "ns/op\t42587780 B/op\t  139658 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50704225,
            "unit": "ns/op\t42625201 B/op\t  139661 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50909192,
            "unit": "ns/op\t42635316 B/op\t  139672 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50928337,
            "unit": "ns/op\t42500923 B/op\t  139583 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41774995,
            "unit": "ns/op\t41676563 B/op\t  156317 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41698269,
            "unit": "ns/op\t41676072 B/op\t  156277 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41587490,
            "unit": "ns/op\t41712984 B/op\t  156335 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41448662,
            "unit": "ns/op\t41744281 B/op\t  156361 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41851728,
            "unit": "ns/op\t41716182 B/op\t  156337 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41948403,
            "unit": "ns/op\t41715392 B/op\t  156329 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41732235,
            "unit": "ns/op\t41723314 B/op\t  156338 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41780633,
            "unit": "ns/op\t41697354 B/op\t  156311 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41382887,
            "unit": "ns/op\t41680044 B/op\t  156306 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41925708,
            "unit": "ns/op\t41741242 B/op\t  156368 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113676418,
            "unit": "ns/op\t61228651 B/op\t  715222 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 112314210,
            "unit": "ns/op\t61101147 B/op\t  715238 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 112436476,
            "unit": "ns/op\t60999852 B/op\t  715205 allocs/op",
            "extra": "10 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114592451,
            "unit": "ns/op\t61235704 B/op\t  715223 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 112914878,
            "unit": "ns/op\t61319842 B/op\t  715281 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "639c77dbee010cd368589f20e700cb7161d252d9",
          "message": "Enable NaN tests (#155)\n\n* Enable NaN tests\r\n\r\nThis commit enables NaN tests by using the latest version of efficientgo/core.\r\n\r\n* Auto-detect NaN values",
          "timestamp": "2023-01-21T07:50:36+01:00",
          "tree_id": "4e6652fc0ee41cd0f71a8b42e46596432827a81b",
          "url": "https://github.com/thanos-community/promql-engine/commit/639c77dbee010cd368589f20e700cb7161d252d9"
        },
        "date": 1674284002177,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22238691,
            "unit": "ns/op\t38819144 B/op\t  131536 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22568980,
            "unit": "ns/op\t38803098 B/op\t  131536 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22392894,
            "unit": "ns/op\t38847967 B/op\t  131539 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22658976,
            "unit": "ns/op\t38852180 B/op\t  131557 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22641582,
            "unit": "ns/op\t38833433 B/op\t  131538 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12161027,
            "unit": "ns/op\t10985403 B/op\t  126100 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12076373,
            "unit": "ns/op\t10941179 B/op\t  126051 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12117249,
            "unit": "ns/op\t11006643 B/op\t  126120 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12159667,
            "unit": "ns/op\t10953289 B/op\t  126069 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12138110,
            "unit": "ns/op\t10978990 B/op\t  126103 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17145259,
            "unit": "ns/op\t24384080 B/op\t  211911 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17197677,
            "unit": "ns/op\t24342664 B/op\t  211893 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17203494,
            "unit": "ns/op\t24377671 B/op\t  211906 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17197825,
            "unit": "ns/op\t24378451 B/op\t  211915 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17177173,
            "unit": "ns/op\t24385145 B/op\t  211916 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12466048,
            "unit": "ns/op\t12278400 B/op\t  133149 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12461111,
            "unit": "ns/op\t12302315 B/op\t  133196 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12474543,
            "unit": "ns/op\t12316341 B/op\t  133192 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12480206,
            "unit": "ns/op\t12353576 B/op\t  133246 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12531936,
            "unit": "ns/op\t12330768 B/op\t  133210 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12440912,
            "unit": "ns/op\t12356399 B/op\t  132848 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12501804,
            "unit": "ns/op\t12342185 B/op\t  132820 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12441771,
            "unit": "ns/op\t12319273 B/op\t  132806 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12447223,
            "unit": "ns/op\t12277963 B/op\t  132689 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12494356,
            "unit": "ns/op\t12323708 B/op\t  132792 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28547245,
            "unit": "ns/op\t41737576 B/op\t  155378 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28468905,
            "unit": "ns/op\t41746909 B/op\t  155398 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28738657,
            "unit": "ns/op\t41772556 B/op\t  155396 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28763768,
            "unit": "ns/op\t41744206 B/op\t  155409 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28567865,
            "unit": "ns/op\t41733743 B/op\t  155371 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24491936,
            "unit": "ns/op\t13710242 B/op\t  149958 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24677549,
            "unit": "ns/op\t13627852 B/op\t  149850 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24320947,
            "unit": "ns/op\t13608154 B/op\t  149830 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24574726,
            "unit": "ns/op\t13654595 B/op\t  149871 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24696385,
            "unit": "ns/op\t13585101 B/op\t  149805 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27886451,
            "unit": "ns/op\t26158838 B/op\t  235257 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27799458,
            "unit": "ns/op\t26235118 B/op\t  235300 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27895939,
            "unit": "ns/op\t26285492 B/op\t  235334 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27806363,
            "unit": "ns/op\t26193580 B/op\t  235272 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27894463,
            "unit": "ns/op\t26289565 B/op\t  235340 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26537352,
            "unit": "ns/op\t42587347 B/op\t  619469 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26510564,
            "unit": "ns/op\t42578155 B/op\t  619485 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26491150,
            "unit": "ns/op\t42592880 B/op\t  619513 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26404034,
            "unit": "ns/op\t42616597 B/op\t  619510 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26312018,
            "unit": "ns/op\t42594284 B/op\t  619495 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17969539,
            "unit": "ns/op\t22984946 B/op\t  108556 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17975879,
            "unit": "ns/op\t22992287 B/op\t  108571 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17961283,
            "unit": "ns/op\t22994448 B/op\t  108565 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17984251,
            "unit": "ns/op\t22990864 B/op\t  108560 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 18042012,
            "unit": "ns/op\t22992324 B/op\t  108570 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40230322,
            "unit": "ns/op\t57931227 B/op\t  202555 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40000918,
            "unit": "ns/op\t57938980 B/op\t  202600 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39723060,
            "unit": "ns/op\t57950202 B/op\t  202564 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 40128942,
            "unit": "ns/op\t57941321 B/op\t  202580 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39863530,
            "unit": "ns/op\t57939040 B/op\t  202583 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33027246,
            "unit": "ns/op\t42416051 B/op\t  135865 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32672311,
            "unit": "ns/op\t42386775 B/op\t  135848 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32467410,
            "unit": "ns/op\t42456689 B/op\t  135888 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32380954,
            "unit": "ns/op\t42506390 B/op\t  135915 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32486250,
            "unit": "ns/op\t42443947 B/op\t  135887 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24523547,
            "unit": "ns/op\t40126706 B/op\t  143829 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24511512,
            "unit": "ns/op\t40116350 B/op\t  143830 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24458039,
            "unit": "ns/op\t40095729 B/op\t  143806 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24343982,
            "unit": "ns/op\t40097777 B/op\t  143785 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24228658,
            "unit": "ns/op\t40125398 B/op\t  143830 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33093271,
            "unit": "ns/op\t42047957 B/op\t  132852 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32471466,
            "unit": "ns/op\t42052788 B/op\t  132869 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32899942,
            "unit": "ns/op\t42134869 B/op\t  132895 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32686278,
            "unit": "ns/op\t42082986 B/op\t  132888 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33035012,
            "unit": "ns/op\t42024999 B/op\t  132838 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22221793,
            "unit": "ns/op\t37053009 B/op\t  102577 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22286442,
            "unit": "ns/op\t37081571 B/op\t  102610 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22226003,
            "unit": "ns/op\t37089049 B/op\t  102604 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 21966519,
            "unit": "ns/op\t37100839 B/op\t  102612 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22256955,
            "unit": "ns/op\t37092132 B/op\t  102610 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17108030,
            "unit": "ns/op\t33224423 B/op\t   74864 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17516326,
            "unit": "ns/op\t33226670 B/op\t   74866 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17647932,
            "unit": "ns/op\t33226383 B/op\t   74864 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17131383,
            "unit": "ns/op\t33223881 B/op\t   74865 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16942603,
            "unit": "ns/op\t33224868 B/op\t   74865 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17139443,
            "unit": "ns/op\t33046412 B/op\t   68867 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16500744,
            "unit": "ns/op\t33045479 B/op\t   68866 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16924855,
            "unit": "ns/op\t33046232 B/op\t   68866 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17304670,
            "unit": "ns/op\t33046323 B/op\t   68867 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17478987,
            "unit": "ns/op\t33044448 B/op\t   68866 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40702169,
            "unit": "ns/op\t39181347 B/op\t  135296 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40345444,
            "unit": "ns/op\t39196381 B/op\t  135330 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40361384,
            "unit": "ns/op\t39229657 B/op\t  135340 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40507114,
            "unit": "ns/op\t39263535 B/op\t  135384 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40138705,
            "unit": "ns/op\t39229836 B/op\t  135376 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37308775,
            "unit": "ns/op\t39168532 B/op\t  134916 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37463812,
            "unit": "ns/op\t39141175 B/op\t  134890 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37532712,
            "unit": "ns/op\t39146972 B/op\t  134908 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37369510,
            "unit": "ns/op\t39157763 B/op\t  134913 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37494411,
            "unit": "ns/op\t39157000 B/op\t  134890 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51715930,
            "unit": "ns/op\t42532565 B/op\t  139653 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51211600,
            "unit": "ns/op\t42700912 B/op\t  139770 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50902231,
            "unit": "ns/op\t42581987 B/op\t  139707 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51324923,
            "unit": "ns/op\t42541562 B/op\t  139662 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50366240,
            "unit": "ns/op\t42544691 B/op\t  139665 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42389422,
            "unit": "ns/op\t41692153 B/op\t  156354 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41708366,
            "unit": "ns/op\t41705948 B/op\t  156347 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41503516,
            "unit": "ns/op\t41728535 B/op\t  156390 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42116594,
            "unit": "ns/op\t41698358 B/op\t  156339 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41492736,
            "unit": "ns/op\t41742591 B/op\t  156407 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42146313,
            "unit": "ns/op\t41709003 B/op\t  156411 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41797383,
            "unit": "ns/op\t41692808 B/op\t  156347 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41953030,
            "unit": "ns/op\t41741565 B/op\t  156397 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42119950,
            "unit": "ns/op\t41710064 B/op\t  156400 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42001096,
            "unit": "ns/op\t41742298 B/op\t  156431 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114745794,
            "unit": "ns/op\t61347000 B/op\t  715296 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114315623,
            "unit": "ns/op\t61124164 B/op\t  715234 allocs/op",
            "extra": "10 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114066335,
            "unit": "ns/op\t61154012 B/op\t  715211 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113366854,
            "unit": "ns/op\t61298800 B/op\t  715265 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113620156,
            "unit": "ns/op\t61400225 B/op\t  715321 allocs/op",
            "extra": "10 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "970995648ccfcc3000a00041ce7a6948bd413bd9",
          "message": "Fix instant query result sort for topk/bottomk (#156)\n\nThere is an implicit expectation that topk and bottomk instant queries\r\nshould return values in sorted order. When a grouping is used for these\r\naggregations, results should be sorted within individual groups.\r\n\r\nThe prometheus engine does not sort the groups themselves, which makes\r\nit hard to do exact comparisons. However, we can now compare results\r\ndirectly for these aggregations when no grouping is used.",
          "timestamp": "2023-01-21T09:22:21+01:00",
          "tree_id": "a5e7b6a8544a9cd8223998b0b6eefd9e027e447f",
          "url": "https://github.com/thanos-community/promql-engine/commit/970995648ccfcc3000a00041ce7a6948bd413bd9"
        },
        "date": 1674289511524,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 21849690,
            "unit": "ns/op\t38783849 B/op\t  131432 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22019157,
            "unit": "ns/op\t38822028 B/op\t  131445 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22258253,
            "unit": "ns/op\t38778873 B/op\t  131410 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22573244,
            "unit": "ns/op\t38805014 B/op\t  131437 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22469922,
            "unit": "ns/op\t38814950 B/op\t  131431 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12081589,
            "unit": "ns/op\t10964227 B/op\t  125997 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12076864,
            "unit": "ns/op\t10977441 B/op\t  126013 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12091599,
            "unit": "ns/op\t10969828 B/op\t  125999 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12088201,
            "unit": "ns/op\t10936522 B/op\t  125966 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12127825,
            "unit": "ns/op\t11020185 B/op\t  126035 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17008400,
            "unit": "ns/op\t24357739 B/op\t  211826 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17076840,
            "unit": "ns/op\t24383192 B/op\t  211810 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17146791,
            "unit": "ns/op\t24363462 B/op\t  211817 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17120865,
            "unit": "ns/op\t24380532 B/op\t  211813 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17147777,
            "unit": "ns/op\t24384317 B/op\t  211834 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12398715,
            "unit": "ns/op\t12395761 B/op\t  138023 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12485342,
            "unit": "ns/op\t12432822 B/op\t  138157 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12467064,
            "unit": "ns/op\t12442502 B/op\t  138082 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12526566,
            "unit": "ns/op\t12434190 B/op\t  138056 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12425464,
            "unit": "ns/op\t12438619 B/op\t  138075 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12389395,
            "unit": "ns/op\t12242495 B/op\t  131675 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12370805,
            "unit": "ns/op\t12300296 B/op\t  131767 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12386138,
            "unit": "ns/op\t12288749 B/op\t  131731 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12383510,
            "unit": "ns/op\t12275591 B/op\t  131719 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12437079,
            "unit": "ns/op\t12291381 B/op\t  131758 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28368697,
            "unit": "ns/op\t41716026 B/op\t  155270 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28671032,
            "unit": "ns/op\t41742201 B/op\t  155317 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28640903,
            "unit": "ns/op\t41737999 B/op\t  155300 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28563207,
            "unit": "ns/op\t41748946 B/op\t  155323 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28706948,
            "unit": "ns/op\t41742261 B/op\t  155286 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24411453,
            "unit": "ns/op\t13693880 B/op\t  149818 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24358315,
            "unit": "ns/op\t13687799 B/op\t  149840 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24473619,
            "unit": "ns/op\t13623921 B/op\t  149778 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24463355,
            "unit": "ns/op\t13646235 B/op\t  149786 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24504960,
            "unit": "ns/op\t13607751 B/op\t  149746 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27676758,
            "unit": "ns/op\t26315451 B/op\t  235285 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27738942,
            "unit": "ns/op\t26182809 B/op\t  235197 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27775846,
            "unit": "ns/op\t26250280 B/op\t  235232 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27700746,
            "unit": "ns/op\t26243045 B/op\t  235223 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27825404,
            "unit": "ns/op\t26189905 B/op\t  235172 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26278104,
            "unit": "ns/op\t42560320 B/op\t  619281 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26243979,
            "unit": "ns/op\t42574754 B/op\t  619284 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26123820,
            "unit": "ns/op\t42537613 B/op\t  619290 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26108934,
            "unit": "ns/op\t42559340 B/op\t  619304 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26267815,
            "unit": "ns/op\t42588567 B/op\t  619277 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17950516,
            "unit": "ns/op\t22995864 B/op\t  108774 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17849459,
            "unit": "ns/op\t22996961 B/op\t  108749 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17806158,
            "unit": "ns/op\t22990995 B/op\t  108748 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17964748,
            "unit": "ns/op\t22993630 B/op\t  108770 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 17935812,
            "unit": "ns/op\t22998919 B/op\t  108772 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39306300,
            "unit": "ns/op\t57903037 B/op\t  202578 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39350053,
            "unit": "ns/op\t57946772 B/op\t  202628 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39716059,
            "unit": "ns/op\t57940840 B/op\t  202652 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39452316,
            "unit": "ns/op\t57924088 B/op\t  202595 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 39665331,
            "unit": "ns/op\t57936058 B/op\t  202645 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32796263,
            "unit": "ns/op\t42425072 B/op\t  135779 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32441794,
            "unit": "ns/op\t42432312 B/op\t  135785 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32458019,
            "unit": "ns/op\t42424095 B/op\t  135785 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32165549,
            "unit": "ns/op\t42440518 B/op\t  135789 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32397954,
            "unit": "ns/op\t42476906 B/op\t  135810 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24564579,
            "unit": "ns/op\t40080435 B/op\t  143705 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24410203,
            "unit": "ns/op\t40084184 B/op\t  143706 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24259229,
            "unit": "ns/op\t40111673 B/op\t  143737 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24377431,
            "unit": "ns/op\t40120522 B/op\t  143737 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24506987,
            "unit": "ns/op\t40121702 B/op\t  143737 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32345325,
            "unit": "ns/op\t42137615 B/op\t  132822 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32326909,
            "unit": "ns/op\t42041769 B/op\t  132756 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32744525,
            "unit": "ns/op\t42136367 B/op\t  132806 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32443938,
            "unit": "ns/op\t41991941 B/op\t  132733 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 32422986,
            "unit": "ns/op\t42048105 B/op\t  132760 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22090615,
            "unit": "ns/op\t37092122 B/op\t  102562 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22004087,
            "unit": "ns/op\t37065911 B/op\t  102530 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 21830740,
            "unit": "ns/op\t37067237 B/op\t  102523 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22193317,
            "unit": "ns/op\t37074999 B/op\t  102556 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22119987,
            "unit": "ns/op\t37096347 B/op\t  102542 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17641443,
            "unit": "ns/op\t33222157 B/op\t   74835 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17314898,
            "unit": "ns/op\t33224012 B/op\t   74834 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17077125,
            "unit": "ns/op\t33221275 B/op\t   74834 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17200642,
            "unit": "ns/op\t33219791 B/op\t   74834 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16945154,
            "unit": "ns/op\t33225111 B/op\t   74835 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16914345,
            "unit": "ns/op\t33046412 B/op\t   68837 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17162089,
            "unit": "ns/op\t33046382 B/op\t   68838 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17330564,
            "unit": "ns/op\t33049558 B/op\t   68838 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16294712,
            "unit": "ns/op\t33041229 B/op\t   68835 allocs/op",
            "extra": "86 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17035102,
            "unit": "ns/op\t33046424 B/op\t   68836 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40730403,
            "unit": "ns/op\t39202343 B/op\t  135234 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 39602148,
            "unit": "ns/op\t39209924 B/op\t  135258 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40422968,
            "unit": "ns/op\t39185790 B/op\t  135214 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40210484,
            "unit": "ns/op\t39192084 B/op\t  135229 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40163826,
            "unit": "ns/op\t39185806 B/op\t  135221 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37593126,
            "unit": "ns/op\t39128332 B/op\t  134800 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37416766,
            "unit": "ns/op\t39154700 B/op\t  134838 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37627911,
            "unit": "ns/op\t39153741 B/op\t  134801 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37225917,
            "unit": "ns/op\t39152458 B/op\t  134824 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 37362360,
            "unit": "ns/op\t39131421 B/op\t  134794 allocs/op",
            "extra": "33 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50954555,
            "unit": "ns/op\t42641029 B/op\t  139638 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50988393,
            "unit": "ns/op\t42586458 B/op\t  139592 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 49976555,
            "unit": "ns/op\t42683167 B/op\t  139657 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50614071,
            "unit": "ns/op\t42564534 B/op\t  139578 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50850359,
            "unit": "ns/op\t42600346 B/op\t  139607 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41420117,
            "unit": "ns/op\t41708239 B/op\t  156292 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41808884,
            "unit": "ns/op\t41728881 B/op\t  156300 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41810183,
            "unit": "ns/op\t41712323 B/op\t  156275 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41837223,
            "unit": "ns/op\t41703695 B/op\t  156273 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 41283328,
            "unit": "ns/op\t41739934 B/op\t  156321 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41550637,
            "unit": "ns/op\t41731847 B/op\t  156323 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41599969,
            "unit": "ns/op\t41710691 B/op\t  156276 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41757594,
            "unit": "ns/op\t41746167 B/op\t  156323 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41599044,
            "unit": "ns/op\t41688951 B/op\t  156253 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 41709410,
            "unit": "ns/op\t41716402 B/op\t  156304 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 112921939,
            "unit": "ns/op\t61292544 B/op\t  715358 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113218594,
            "unit": "ns/op\t61424577 B/op\t  715331 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113656789,
            "unit": "ns/op\t61304664 B/op\t  715364 allocs/op",
            "extra": "10 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114459386,
            "unit": "ns/op\t61143003 B/op\t  715319 allocs/op",
            "extra": "10 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113145389,
            "unit": "ns/op\t61254792 B/op\t  715352 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ee0f08df788b061999d6dde8873caa5edb55c42e",
          "message": "Implement distributed deduplication (#151)\n\n* Implement distributed deduplication\r\n\r\nThis commit adds support for deduplicating series/samples and applies\r\nit to the distributed query engine. This is done by implementing a\r\nDeduplication operator that removes samples with the same sample IDs from\r\nindividual step vectors. Samples are deduplicated by a\r\nlast-sample-wins strategy.\r\n\r\nIn order to properly prioritize non-stale samples in favor of stale ones,\r\nremote engines are sorted by their max time. This way if samples come\r\nfrom both older and newer engines, the ones from newer engines will be\r\npreserved. Because of thism, stale samples will be deduplicated in favour\r\nof non-stale ones, but will be preserved if non-stale samples are not\r\npresent inside a step vector.\r\n\r\n* Add license header\r\n\r\n* Fix logical tests\r\n\r\n* Remove implicit dependency between logical plan and execution\r\n\r\n* Add dedup comment",
          "timestamp": "2023-01-23T10:46:37+01:00",
          "tree_id": "8f49ba898668251b7ca9ac3fab2e09fae189cda6",
          "url": "https://github.com/thanos-community/promql-engine/commit/ee0f08df788b061999d6dde8873caa5edb55c42e"
        },
        "date": 1674467372683,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22861723,
            "unit": "ns/op\t39390844 B/op\t  131389 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23171885,
            "unit": "ns/op\t39365072 B/op\t  131354 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23317986,
            "unit": "ns/op\t39416400 B/op\t  131418 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23535233,
            "unit": "ns/op\t39412834 B/op\t  131414 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23442384,
            "unit": "ns/op\t39373615 B/op\t  131365 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12290570,
            "unit": "ns/op\t11306978 B/op\t  126003 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12266208,
            "unit": "ns/op\t11296491 B/op\t  125980 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12174239,
            "unit": "ns/op\t11272322 B/op\t  125973 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12211959,
            "unit": "ns/op\t11266066 B/op\t  125973 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12247893,
            "unit": "ns/op\t11295478 B/op\t  125994 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17669567,
            "unit": "ns/op\t24558670 B/op\t  211773 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17699338,
            "unit": "ns/op\t24550672 B/op\t  211774 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17724476,
            "unit": "ns/op\t24548492 B/op\t  211787 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17816901,
            "unit": "ns/op\t24566103 B/op\t  211780 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17670153,
            "unit": "ns/op\t24573524 B/op\t  211795 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13018861,
            "unit": "ns/op\t12800779 B/op\t  135315 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12974092,
            "unit": "ns/op\t12816285 B/op\t  135341 allocs/op",
            "extra": "86 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12952599,
            "unit": "ns/op\t12798541 B/op\t  135324 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12873748,
            "unit": "ns/op\t12791988 B/op\t  135309 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12915388,
            "unit": "ns/op\t12778731 B/op\t  135305 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12977534,
            "unit": "ns/op\t12818950 B/op\t  135133 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12929918,
            "unit": "ns/op\t12783067 B/op\t  135085 allocs/op",
            "extra": "87 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12988700,
            "unit": "ns/op\t12831109 B/op\t  135144 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12842913,
            "unit": "ns/op\t12747303 B/op\t  135052 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12925604,
            "unit": "ns/op\t12818717 B/op\t  135111 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28638658,
            "unit": "ns/op\t42200507 B/op\t  155252 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28913723,
            "unit": "ns/op\t42242047 B/op\t  155308 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 29050201,
            "unit": "ns/op\t42218823 B/op\t  155271 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28849237,
            "unit": "ns/op\t42208286 B/op\t  155270 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28762509,
            "unit": "ns/op\t42218076 B/op\t  155285 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24774716,
            "unit": "ns/op\t13952667 B/op\t  149797 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24787845,
            "unit": "ns/op\t13934807 B/op\t  149775 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24546772,
            "unit": "ns/op\t13925707 B/op\t  149764 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24656533,
            "unit": "ns/op\t14002436 B/op\t  149835 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24717824,
            "unit": "ns/op\t13922837 B/op\t  149761 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28035694,
            "unit": "ns/op\t26363604 B/op\t  235120 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28147608,
            "unit": "ns/op\t26376227 B/op\t  235153 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28175029,
            "unit": "ns/op\t26340790 B/op\t  235104 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28145703,
            "unit": "ns/op\t26453562 B/op\t  235187 allocs/op",
            "extra": "40 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28087329,
            "unit": "ns/op\t26346133 B/op\t  235114 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27138665,
            "unit": "ns/op\t42867306 B/op\t  619188 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26976805,
            "unit": "ns/op\t42881614 B/op\t  619217 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26833938,
            "unit": "ns/op\t42833377 B/op\t  619145 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26916004,
            "unit": "ns/op\t42857472 B/op\t  619205 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26958822,
            "unit": "ns/op\t42836491 B/op\t  619152 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19588584,
            "unit": "ns/op\t23123447 B/op\t  108985 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19674003,
            "unit": "ns/op\t23110793 B/op\t  108947 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19640390,
            "unit": "ns/op\t23115951 B/op\t  108954 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19437184,
            "unit": "ns/op\t23111694 B/op\t  108946 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19637225,
            "unit": "ns/op\t23124592 B/op\t  108969 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43737853,
            "unit": "ns/op\t58093863 B/op\t  202735 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43164025,
            "unit": "ns/op\t58173517 B/op\t  202852 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42705442,
            "unit": "ns/op\t58145252 B/op\t  202817 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43394608,
            "unit": "ns/op\t58127318 B/op\t  202785 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43588558,
            "unit": "ns/op\t58063054 B/op\t  202718 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 34062053,
            "unit": "ns/op\t42734094 B/op\t  135631 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33703733,
            "unit": "ns/op\t42678186 B/op\t  135601 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33419140,
            "unit": "ns/op\t42776311 B/op\t  135654 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33572180,
            "unit": "ns/op\t42710508 B/op\t  135618 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33476919,
            "unit": "ns/op\t42657473 B/op\t  135587 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24995858,
            "unit": "ns/op\t40797978 B/op\t  143864 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24546396,
            "unit": "ns/op\t40770809 B/op\t  143848 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24589701,
            "unit": "ns/op\t40824184 B/op\t  143905 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24782302,
            "unit": "ns/op\t40791833 B/op\t  143861 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24897031,
            "unit": "ns/op\t40795171 B/op\t  143862 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33885101,
            "unit": "ns/op\t42378652 B/op\t  132623 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34078911,
            "unit": "ns/op\t42353324 B/op\t  132615 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33739885,
            "unit": "ns/op\t42362839 B/op\t  132615 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34010961,
            "unit": "ns/op\t42465647 B/op\t  132682 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33711144,
            "unit": "ns/op\t42398048 B/op\t  132643 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23009810,
            "unit": "ns/op\t37648483 B/op\t  102573 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22941628,
            "unit": "ns/op\t37611062 B/op\t  102534 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22719797,
            "unit": "ns/op\t37613867 B/op\t  102535 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22913877,
            "unit": "ns/op\t37647416 B/op\t  102571 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22939870,
            "unit": "ns/op\t37603935 B/op\t  102520 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17494057,
            "unit": "ns/op\t33218103 B/op\t   74773 allocs/op",
            "extra": "78 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17116381,
            "unit": "ns/op\t33218785 B/op\t   74773 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17461787,
            "unit": "ns/op\t33220262 B/op\t   74773 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16950616,
            "unit": "ns/op\t33217594 B/op\t   74774 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16958834,
            "unit": "ns/op\t33220297 B/op\t   74775 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16967060,
            "unit": "ns/op\t33043844 B/op\t   68775 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16906668,
            "unit": "ns/op\t33040784 B/op\t   68776 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17010363,
            "unit": "ns/op\t33043439 B/op\t   68777 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16984602,
            "unit": "ns/op\t33042538 B/op\t   68776 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16916531,
            "unit": "ns/op\t33038391 B/op\t   68775 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41436686,
            "unit": "ns/op\t39687535 B/op\t  135086 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41785377,
            "unit": "ns/op\t39684218 B/op\t  135080 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41418139,
            "unit": "ns/op\t39711752 B/op\t  135130 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41190365,
            "unit": "ns/op\t39700727 B/op\t  135109 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41240425,
            "unit": "ns/op\t39729957 B/op\t  135150 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38515874,
            "unit": "ns/op\t39680439 B/op\t  134704 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38251337,
            "unit": "ns/op\t39711281 B/op\t  134743 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39119657,
            "unit": "ns/op\t39678196 B/op\t  134702 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38824230,
            "unit": "ns/op\t39642384 B/op\t  134650 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39148114,
            "unit": "ns/op\t39666447 B/op\t  134690 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51947665,
            "unit": "ns/op\t42933305 B/op\t  139479 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52143237,
            "unit": "ns/op\t42824526 B/op\t  139394 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52513813,
            "unit": "ns/op\t43056353 B/op\t  139564 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52553842,
            "unit": "ns/op\t43040057 B/op\t  139554 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52436736,
            "unit": "ns/op\t42941634 B/op\t  139478 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42779800,
            "unit": "ns/op\t42187677 B/op\t  156033 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42767998,
            "unit": "ns/op\t42205465 B/op\t  156057 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43059306,
            "unit": "ns/op\t42219982 B/op\t  156066 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43195366,
            "unit": "ns/op\t42186618 B/op\t  156028 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42722906,
            "unit": "ns/op\t42212348 B/op\t  156068 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43338080,
            "unit": "ns/op\t42222811 B/op\t  156071 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43378089,
            "unit": "ns/op\t42217660 B/op\t  156068 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42915659,
            "unit": "ns/op\t42196216 B/op\t  156041 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43379872,
            "unit": "ns/op\t42230584 B/op\t  156087 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43244566,
            "unit": "ns/op\t42226896 B/op\t  156095 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115135555,
            "unit": "ns/op\t61945896 B/op\t  715073 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115592100,
            "unit": "ns/op\t62656609 B/op\t  715360 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115946342,
            "unit": "ns/op\t61993933 B/op\t  715089 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114366520,
            "unit": "ns/op\t62049700 B/op\t  715119 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115807799,
            "unit": "ns/op\t62509560 B/op\t  715304 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "filip.petkovsky@gmail.com",
            "name": "Filip Petkovski",
            "username": "fpetkovski"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0ad3f3b2e4b46115f6609c4e03af3ad30834c444",
          "message": "Fix distributing an aggregation inside a function (#158)\n\nThe distribute optimizer is passing a reference copy when traversing\r\nfunction arguments. Even when the copy is modified, the original node\r\nremains the same. This causes queries with histogram_quantile to get\r\nexecuted in memory instead of getting distributed to remote engines.",
          "timestamp": "2023-01-23T15:36:05+01:00",
          "tree_id": "c931c387fb8d1b4db5ac021bc4f25f2d49c4d929",
          "url": "https://github.com/thanos-community/promql-engine/commit/0ad3f3b2e4b46115f6609c4e03af3ad30834c444"
        },
        "date": 1674484738205,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22511276,
            "unit": "ns/op\t39362177 B/op\t  131552 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22582228,
            "unit": "ns/op\t39384767 B/op\t  131570 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22476074,
            "unit": "ns/op\t39349696 B/op\t  131531 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22421828,
            "unit": "ns/op\t39386356 B/op\t  131568 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23281909,
            "unit": "ns/op\t39366573 B/op\t  131551 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12157580,
            "unit": "ns/op\t11284732 B/op\t  126178 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12136270,
            "unit": "ns/op\t11290503 B/op\t  126177 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12206739,
            "unit": "ns/op\t11298618 B/op\t  126188 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12171602,
            "unit": "ns/op\t11314564 B/op\t  126200 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12197439,
            "unit": "ns/op\t11271138 B/op\t  126159 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17535277,
            "unit": "ns/op\t24569648 B/op\t  211985 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17578392,
            "unit": "ns/op\t24567968 B/op\t  211982 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17793540,
            "unit": "ns/op\t24559496 B/op\t  211978 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17689510,
            "unit": "ns/op\t24575616 B/op\t  211973 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17692934,
            "unit": "ns/op\t24584126 B/op\t  211998 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12729701,
            "unit": "ns/op\t12695885 B/op\t  130834 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12749992,
            "unit": "ns/op\t12711103 B/op\t  130845 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12781405,
            "unit": "ns/op\t12672785 B/op\t  130817 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12688197,
            "unit": "ns/op\t12669660 B/op\t  130801 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12791835,
            "unit": "ns/op\t12726726 B/op\t  130871 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12822672,
            "unit": "ns/op\t12835558 B/op\t  135720 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12925346,
            "unit": "ns/op\t12835093 B/op\t  135727 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12856287,
            "unit": "ns/op\t12804220 B/op\t  135707 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12806093,
            "unit": "ns/op\t12771717 B/op\t  135663 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12838846,
            "unit": "ns/op\t12775816 B/op\t  135667 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28590103,
            "unit": "ns/op\t42238806 B/op\t  155499 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28601463,
            "unit": "ns/op\t42195874 B/op\t  155454 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28625553,
            "unit": "ns/op\t42228337 B/op\t  155491 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28571638,
            "unit": "ns/op\t42183607 B/op\t  155435 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28605973,
            "unit": "ns/op\t42198543 B/op\t  155447 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24656381,
            "unit": "ns/op\t13877889 B/op\t  149928 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24493658,
            "unit": "ns/op\t13902768 B/op\t  149958 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24621887,
            "unit": "ns/op\t13929706 B/op\t  149963 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24460002,
            "unit": "ns/op\t13982255 B/op\t  150029 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24805032,
            "unit": "ns/op\t13885634 B/op\t  149931 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27770178,
            "unit": "ns/op\t26378166 B/op\t  235326 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28134650,
            "unit": "ns/op\t26311296 B/op\t  235274 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28094141,
            "unit": "ns/op\t26291670 B/op\t  235271 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28020486,
            "unit": "ns/op\t26250772 B/op\t  235240 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28073490,
            "unit": "ns/op\t26347497 B/op\t  235296 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26857506,
            "unit": "ns/op\t42882494 B/op\t  619584 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26833381,
            "unit": "ns/op\t42877693 B/op\t  619610 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26787176,
            "unit": "ns/op\t42872315 B/op\t  619596 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26924023,
            "unit": "ns/op\t42898845 B/op\t  619603 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26851108,
            "unit": "ns/op\t42864762 B/op\t  619610 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19555283,
            "unit": "ns/op\t23106054 B/op\t  108664 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19520028,
            "unit": "ns/op\t23093107 B/op\t  108630 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19425755,
            "unit": "ns/op\t23105088 B/op\t  108678 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19525867,
            "unit": "ns/op\t23101838 B/op\t  108660 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19363679,
            "unit": "ns/op\t23110485 B/op\t  108676 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42832990,
            "unit": "ns/op\t58154407 B/op\t  202794 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42971236,
            "unit": "ns/op\t58088312 B/op\t  202714 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43113710,
            "unit": "ns/op\t58083517 B/op\t  202689 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42566673,
            "unit": "ns/op\t58110632 B/op\t  202738 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42397896,
            "unit": "ns/op\t58136136 B/op\t  202773 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33816778,
            "unit": "ns/op\t42666907 B/op\t  135782 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33411785,
            "unit": "ns/op\t42623044 B/op\t  135768 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33604364,
            "unit": "ns/op\t42629702 B/op\t  135772 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33620190,
            "unit": "ns/op\t42704463 B/op\t  135815 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33132699,
            "unit": "ns/op\t42834941 B/op\t  135887 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24626935,
            "unit": "ns/op\t40759048 B/op\t  144001 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24731126,
            "unit": "ns/op\t40790629 B/op\t  144059 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24633975,
            "unit": "ns/op\t40781381 B/op\t  144042 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24671379,
            "unit": "ns/op\t40801846 B/op\t  144078 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24665934,
            "unit": "ns/op\t40778576 B/op\t  144054 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33350991,
            "unit": "ns/op\t42406480 B/op\t  132847 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33724504,
            "unit": "ns/op\t42343818 B/op\t  132810 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34047774,
            "unit": "ns/op\t42350237 B/op\t  132807 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33453524,
            "unit": "ns/op\t42280099 B/op\t  132772 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33678437,
            "unit": "ns/op\t42367207 B/op\t  132821 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22752050,
            "unit": "ns/op\t37649971 B/op\t  102704 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22727731,
            "unit": "ns/op\t37591034 B/op\t  102643 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22594310,
            "unit": "ns/op\t37636239 B/op\t  102684 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22665045,
            "unit": "ns/op\t37650729 B/op\t  102699 allocs/op",
            "extra": "60 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22648271,
            "unit": "ns/op\t37611332 B/op\t  102657 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 18092190,
            "unit": "ns/op\t33218785 B/op\t   74841 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17286004,
            "unit": "ns/op\t33221016 B/op\t   74840 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17407729,
            "unit": "ns/op\t33220990 B/op\t   74839 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17849884,
            "unit": "ns/op\t33221318 B/op\t   74840 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17451191,
            "unit": "ns/op\t33221537 B/op\t   74840 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17552973,
            "unit": "ns/op\t33043155 B/op\t   68842 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17354290,
            "unit": "ns/op\t33040738 B/op\t   68841 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17328158,
            "unit": "ns/op\t33041008 B/op\t   68841 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17084967,
            "unit": "ns/op\t33043769 B/op\t   68841 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16945822,
            "unit": "ns/op\t33038501 B/op\t   68840 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41678737,
            "unit": "ns/op\t39677517 B/op\t  135269 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41164073,
            "unit": "ns/op\t39736208 B/op\t  135346 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41922863,
            "unit": "ns/op\t39695714 B/op\t  135288 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41990823,
            "unit": "ns/op\t39660520 B/op\t  135243 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41657158,
            "unit": "ns/op\t39701625 B/op\t  135304 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38952664,
            "unit": "ns/op\t39667388 B/op\t  134882 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38819537,
            "unit": "ns/op\t39663086 B/op\t  134870 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38820934,
            "unit": "ns/op\t39682237 B/op\t  134902 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38781543,
            "unit": "ns/op\t39676406 B/op\t  134899 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38856618,
            "unit": "ns/op\t39690167 B/op\t  134909 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52090458,
            "unit": "ns/op\t42986007 B/op\t  139704 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52801019,
            "unit": "ns/op\t42955511 B/op\t  139676 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52045005,
            "unit": "ns/op\t43082952 B/op\t  139781 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51516669,
            "unit": "ns/op\t43019368 B/op\t  139735 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52890931,
            "unit": "ns/op\t42938069 B/op\t  139682 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42877555,
            "unit": "ns/op\t42226774 B/op\t  156272 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42877149,
            "unit": "ns/op\t42236601 B/op\t  156293 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42745934,
            "unit": "ns/op\t42233719 B/op\t  156282 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43035210,
            "unit": "ns/op\t42185047 B/op\t  156222 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43057422,
            "unit": "ns/op\t42223431 B/op\t  156273 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43370259,
            "unit": "ns/op\t42213200 B/op\t  156258 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42439353,
            "unit": "ns/op\t42170838 B/op\t  156198 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43041431,
            "unit": "ns/op\t42222329 B/op\t  156274 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43219845,
            "unit": "ns/op\t42210914 B/op\t  156247 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42554236,
            "unit": "ns/op\t42201840 B/op\t  156237 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113830892,
            "unit": "ns/op\t61996780 B/op\t  715438 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114845557,
            "unit": "ns/op\t62153372 B/op\t  715509 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115044864,
            "unit": "ns/op\t62133855 B/op\t  715475 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114231412,
            "unit": "ns/op\t61788657 B/op\t  715353 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114758937,
            "unit": "ns/op\t62214728 B/op\t  715530 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "121007071+kama910@users.noreply.github.com",
            "name": "Kama Huang",
            "username": "kama910"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9e293186b7e42d365d216df929d274a7bb04d414",
          "message": "add sgn function (#160)\n\n* add sgn function\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>\r\n\r\n* use simplified sgn implementation\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>",
          "timestamp": "2023-01-24T08:04:17+01:00",
          "tree_id": "f5cf4a4ace244e201ac6878cab672e335ec02eba",
          "url": "https://github.com/thanos-community/promql-engine/commit/9e293186b7e42d365d216df929d274a7bb04d414"
        },
        "date": 1674544068883,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22387080,
            "unit": "ns/op\t39345085 B/op\t  131658 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22324420,
            "unit": "ns/op\t39373966 B/op\t  131687 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23358103,
            "unit": "ns/op\t39404680 B/op\t  131718 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23128484,
            "unit": "ns/op\t39406743 B/op\t  131720 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23244532,
            "unit": "ns/op\t39367838 B/op\t  131684 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12387409,
            "unit": "ns/op\t11257622 B/op\t  126270 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12316476,
            "unit": "ns/op\t11278610 B/op\t  126302 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12366631,
            "unit": "ns/op\t11236081 B/op\t  126266 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12332298,
            "unit": "ns/op\t11282737 B/op\t  126305 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12385078,
            "unit": "ns/op\t11266903 B/op\t  126302 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17605709,
            "unit": "ns/op\t24576112 B/op\t  212100 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17541534,
            "unit": "ns/op\t24594522 B/op\t  212122 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17618103,
            "unit": "ns/op\t24572580 B/op\t  212121 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17559101,
            "unit": "ns/op\t24586056 B/op\t  212110 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17513476,
            "unit": "ns/op\t24589543 B/op\t  212112 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12943227,
            "unit": "ns/op\t12774865 B/op\t  135742 allocs/op",
            "extra": "86 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12890792,
            "unit": "ns/op\t12720117 B/op\t  135704 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12866670,
            "unit": "ns/op\t12709395 B/op\t  135701 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12865805,
            "unit": "ns/op\t12761291 B/op\t  135744 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12896234,
            "unit": "ns/op\t12733089 B/op\t  135710 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12631368,
            "unit": "ns/op\t12591960 B/op\t  131146 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12787080,
            "unit": "ns/op\t12579742 B/op\t  131130 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12756402,
            "unit": "ns/op\t12614004 B/op\t  131159 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12785087,
            "unit": "ns/op\t12589863 B/op\t  131129 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12625949,
            "unit": "ns/op\t12604544 B/op\t  131151 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28388854,
            "unit": "ns/op\t42211565 B/op\t  155597 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28276435,
            "unit": "ns/op\t42238067 B/op\t  155629 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28228079,
            "unit": "ns/op\t42210646 B/op\t  155592 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28219765,
            "unit": "ns/op\t42216222 B/op\t  155598 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28459550,
            "unit": "ns/op\t42235692 B/op\t  155633 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24223765,
            "unit": "ns/op\t13923825 B/op\t  150103 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24352983,
            "unit": "ns/op\t13896867 B/op\t  150080 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24228639,
            "unit": "ns/op\t13954907 B/op\t  150124 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24313763,
            "unit": "ns/op\t13876153 B/op\t  150047 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24299050,
            "unit": "ns/op\t13910665 B/op\t  150081 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27623356,
            "unit": "ns/op\t26311176 B/op\t  235413 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27643490,
            "unit": "ns/op\t26480879 B/op\t  235526 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27753705,
            "unit": "ns/op\t26316191 B/op\t  235432 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27580027,
            "unit": "ns/op\t26304733 B/op\t  235433 allocs/op",
            "extra": "58 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27696032,
            "unit": "ns/op\t26328771 B/op\t  235426 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27061135,
            "unit": "ns/op\t42899900 B/op\t  619871 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26860547,
            "unit": "ns/op\t42849373 B/op\t  619850 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26826950,
            "unit": "ns/op\t42919159 B/op\t  619880 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26804961,
            "unit": "ns/op\t42906193 B/op\t  619945 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26698546,
            "unit": "ns/op\t42874097 B/op\t  619867 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19569391,
            "unit": "ns/op\t23101234 B/op\t  108545 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19443171,
            "unit": "ns/op\t23113381 B/op\t  108573 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19458871,
            "unit": "ns/op\t23122742 B/op\t  108606 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19595797,
            "unit": "ns/op\t23097412 B/op\t  108530 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19530739,
            "unit": "ns/op\t23123416 B/op\t  108605 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42684678,
            "unit": "ns/op\t58133821 B/op\t  202760 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42383675,
            "unit": "ns/op\t58151098 B/op\t  202767 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42615839,
            "unit": "ns/op\t58149409 B/op\t  202749 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42275518,
            "unit": "ns/op\t58145202 B/op\t  202766 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42457656,
            "unit": "ns/op\t58137880 B/op\t  202759 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32934665,
            "unit": "ns/op\t42684370 B/op\t  135926 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 32869803,
            "unit": "ns/op\t42645402 B/op\t  135911 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33445997,
            "unit": "ns/op\t42775412 B/op\t  135990 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33035579,
            "unit": "ns/op\t42772830 B/op\t  136001 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33012813,
            "unit": "ns/op\t42662705 B/op\t  135927 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24574437,
            "unit": "ns/op\t40769681 B/op\t  144152 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24417688,
            "unit": "ns/op\t40781581 B/op\t  144175 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24347005,
            "unit": "ns/op\t40765867 B/op\t  144149 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24247844,
            "unit": "ns/op\t40753982 B/op\t  144149 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24601204,
            "unit": "ns/op\t40754211 B/op\t  144160 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33303158,
            "unit": "ns/op\t42321427 B/op\t  132930 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34049727,
            "unit": "ns/op\t42263209 B/op\t  132895 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33113350,
            "unit": "ns/op\t42294606 B/op\t  132909 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33713806,
            "unit": "ns/op\t42377592 B/op\t  132962 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33606742,
            "unit": "ns/op\t42211703 B/op\t  132867 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22696317,
            "unit": "ns/op\t37587033 B/op\t  102714 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22400372,
            "unit": "ns/op\t37592586 B/op\t  102721 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22547836,
            "unit": "ns/op\t37609284 B/op\t  102744 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22586345,
            "unit": "ns/op\t37611377 B/op\t  102747 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22862626,
            "unit": "ns/op\t37628009 B/op\t  102766 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17113858,
            "unit": "ns/op\t33220680 B/op\t   74887 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17402529,
            "unit": "ns/op\t33224711 B/op\t   74888 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17154406,
            "unit": "ns/op\t33224818 B/op\t   74888 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17355315,
            "unit": "ns/op\t33226269 B/op\t   74888 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17369892,
            "unit": "ns/op\t33223321 B/op\t   74886 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17258354,
            "unit": "ns/op\t33045710 B/op\t   68889 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16697734,
            "unit": "ns/op\t33045678 B/op\t   68889 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16853529,
            "unit": "ns/op\t33044151 B/op\t   68891 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16735850,
            "unit": "ns/op\t33046938 B/op\t   68890 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16946228,
            "unit": "ns/op\t33047158 B/op\t   68890 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41798163,
            "unit": "ns/op\t39661352 B/op\t  135377 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41559054,
            "unit": "ns/op\t39680355 B/op\t  135402 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41746831,
            "unit": "ns/op\t39674917 B/op\t  135391 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40983284,
            "unit": "ns/op\t39698386 B/op\t  135424 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41101168,
            "unit": "ns/op\t39715999 B/op\t  135458 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38769295,
            "unit": "ns/op\t39698156 B/op\t  135046 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39095930,
            "unit": "ns/op\t39702071 B/op\t  135056 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39214256,
            "unit": "ns/op\t39652362 B/op\t  134992 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38816460,
            "unit": "ns/op\t39669305 B/op\t  135017 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38740530,
            "unit": "ns/op\t39678653 B/op\t  135026 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51412768,
            "unit": "ns/op\t42881429 B/op\t  139763 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52317721,
            "unit": "ns/op\t42740072 B/op\t  139658 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51320991,
            "unit": "ns/op\t42969216 B/op\t  139832 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51776072,
            "unit": "ns/op\t42834408 B/op\t  139733 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51255165,
            "unit": "ns/op\t42851518 B/op\t  139749 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43094246,
            "unit": "ns/op\t42207465 B/op\t  156372 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42493621,
            "unit": "ns/op\t42220360 B/op\t  156404 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42857736,
            "unit": "ns/op\t42201157 B/op\t  156369 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42838803,
            "unit": "ns/op\t42254749 B/op\t  156443 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42967459,
            "unit": "ns/op\t42200175 B/op\t  156369 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42687066,
            "unit": "ns/op\t42257525 B/op\t  156455 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42622574,
            "unit": "ns/op\t42212012 B/op\t  156379 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42482027,
            "unit": "ns/op\t42204327 B/op\t  156377 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42605918,
            "unit": "ns/op\t42280399 B/op\t  156486 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42390516,
            "unit": "ns/op\t42231576 B/op\t  156420 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113926026,
            "unit": "ns/op\t62170310 B/op\t  715215 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113915158,
            "unit": "ns/op\t62026222 B/op\t  715173 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114920964,
            "unit": "ns/op\t62335707 B/op\t  715280 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115064131,
            "unit": "ns/op\t62393245 B/op\t  715320 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113893524,
            "unit": "ns/op\t62101092 B/op\t  715200 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benye@amazon.com",
            "name": "Ben Ye",
            "username": "yeya24"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d2d4da1b91ae32f878db14f1829192a9b1bdc166",
          "message": "remove fmt.Sprintf() when writing to hash (#164)\n\nSigned-off-by: Ben Ye <benye@amazon.com>",
          "timestamp": "2023-01-28T08:30:46+01:00",
          "tree_id": "eeb192a08b4b5d2fa5e83952dc5a6aa287e93ab7",
          "url": "https://github.com/thanos-community/promql-engine/commit/d2d4da1b91ae32f878db14f1829192a9b1bdc166"
        },
        "date": 1674891216167,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23066590,
            "unit": "ns/op\t39410958 B/op\t  131700 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23242600,
            "unit": "ns/op\t39428962 B/op\t  131718 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23506633,
            "unit": "ns/op\t39379713 B/op\t  131662 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23644924,
            "unit": "ns/op\t39405421 B/op\t  131690 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23684615,
            "unit": "ns/op\t39393781 B/op\t  131678 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12274921,
            "unit": "ns/op\t11319360 B/op\t  126325 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12292983,
            "unit": "ns/op\t11303262 B/op\t  126292 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12300972,
            "unit": "ns/op\t11281074 B/op\t  126272 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12247211,
            "unit": "ns/op\t11301271 B/op\t  126305 allocs/op",
            "extra": "97 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12292061,
            "unit": "ns/op\t11297311 B/op\t  126292 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17818574,
            "unit": "ns/op\t24582045 B/op\t  212076 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17717780,
            "unit": "ns/op\t24584025 B/op\t  212092 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17705416,
            "unit": "ns/op\t24553517 B/op\t  212059 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17649490,
            "unit": "ns/op\t24580659 B/op\t  212080 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17826071,
            "unit": "ns/op\t24579768 B/op\t  212076 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12908838,
            "unit": "ns/op\t12736236 B/op\t  132607 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12924176,
            "unit": "ns/op\t12724652 B/op\t  132578 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12848587,
            "unit": "ns/op\t12676200 B/op\t  132542 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12893665,
            "unit": "ns/op\t12707357 B/op\t  132573 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12866122,
            "unit": "ns/op\t12731377 B/op\t  132597 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12905089,
            "unit": "ns/op\t12746308 B/op\t  133567 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12923221,
            "unit": "ns/op\t12729506 B/op\t  133572 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12935871,
            "unit": "ns/op\t12748421 B/op\t  133576 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12899950,
            "unit": "ns/op\t12741072 B/op\t  133567 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12867379,
            "unit": "ns/op\t12728612 B/op\t  133565 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28582950,
            "unit": "ns/op\t42225739 B/op\t  155577 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28473313,
            "unit": "ns/op\t42176533 B/op\t  155527 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28410429,
            "unit": "ns/op\t42206256 B/op\t  155550 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28486750,
            "unit": "ns/op\t42226387 B/op\t  155582 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28438638,
            "unit": "ns/op\t42229612 B/op\t  155584 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24417714,
            "unit": "ns/op\t13954171 B/op\t  150070 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24477966,
            "unit": "ns/op\t13937510 B/op\t  150076 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24127043,
            "unit": "ns/op\t13950385 B/op\t  150084 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24531648,
            "unit": "ns/op\t13942380 B/op\t  150078 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24549648,
            "unit": "ns/op\t13940908 B/op\t  150062 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27800946,
            "unit": "ns/op\t26372139 B/op\t  235425 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27708393,
            "unit": "ns/op\t26433927 B/op\t  235477 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27745872,
            "unit": "ns/op\t26518615 B/op\t  235515 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27855504,
            "unit": "ns/op\t26357211 B/op\t  235426 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27729359,
            "unit": "ns/op\t26512632 B/op\t  235519 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27166033,
            "unit": "ns/op\t42869846 B/op\t  619762 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27077743,
            "unit": "ns/op\t42892747 B/op\t  619798 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26770499,
            "unit": "ns/op\t42884499 B/op\t  619805 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26820635,
            "unit": "ns/op\t42871950 B/op\t  619807 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26994994,
            "unit": "ns/op\t42874108 B/op\t  619779 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19696468,
            "unit": "ns/op\t23125309 B/op\t  108821 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19585073,
            "unit": "ns/op\t23091363 B/op\t  108746 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19678429,
            "unit": "ns/op\t23107550 B/op\t  108778 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19614971,
            "unit": "ns/op\t23112256 B/op\t  108799 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19666094,
            "unit": "ns/op\t23099413 B/op\t  108751 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43652923,
            "unit": "ns/op\t58072058 B/op\t  202712 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43605064,
            "unit": "ns/op\t58116366 B/op\t  202779 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42976535,
            "unit": "ns/op\t58149126 B/op\t  202830 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43297581,
            "unit": "ns/op\t58112455 B/op\t  202772 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43111379,
            "unit": "ns/op\t58110364 B/op\t  202776 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33627220,
            "unit": "ns/op\t42804471 B/op\t  135973 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33298159,
            "unit": "ns/op\t42785360 B/op\t  135966 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33417136,
            "unit": "ns/op\t42734785 B/op\t  135939 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33427349,
            "unit": "ns/op\t42741998 B/op\t  135934 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33514674,
            "unit": "ns/op\t42768343 B/op\t  135942 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24858633,
            "unit": "ns/op\t40770319 B/op\t  144148 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24706710,
            "unit": "ns/op\t40787192 B/op\t  144164 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24734477,
            "unit": "ns/op\t40756330 B/op\t  144119 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24774074,
            "unit": "ns/op\t40784003 B/op\t  144157 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24798526,
            "unit": "ns/op\t40797219 B/op\t  144158 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33887810,
            "unit": "ns/op\t42335250 B/op\t  132913 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34045535,
            "unit": "ns/op\t42301402 B/op\t  132885 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34063296,
            "unit": "ns/op\t42455461 B/op\t  132970 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33817072,
            "unit": "ns/op\t42298631 B/op\t  132877 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33871148,
            "unit": "ns/op\t42358886 B/op\t  132912 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22921784,
            "unit": "ns/op\t37623133 B/op\t  102729 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22962610,
            "unit": "ns/op\t37588833 B/op\t  102697 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22983233,
            "unit": "ns/op\t37613789 B/op\t  102725 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22712382,
            "unit": "ns/op\t37601596 B/op\t  102707 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22761515,
            "unit": "ns/op\t37621339 B/op\t  102736 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17760787,
            "unit": "ns/op\t33222423 B/op\t   74871 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17220679,
            "unit": "ns/op\t33221586 B/op\t   74871 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17182893,
            "unit": "ns/op\t33225298 B/op\t   74871 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17017217,
            "unit": "ns/op\t33223952 B/op\t   74872 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17256232,
            "unit": "ns/op\t33217863 B/op\t   74870 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17342737,
            "unit": "ns/op\t33045115 B/op\t   68874 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17349065,
            "unit": "ns/op\t33045337 B/op\t   68873 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16981326,
            "unit": "ns/op\t33042889 B/op\t   68873 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16911425,
            "unit": "ns/op\t33043379 B/op\t   68874 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17003121,
            "unit": "ns/op\t33048779 B/op\t   68875 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41643326,
            "unit": "ns/op\t39751925 B/op\t  135475 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41246093,
            "unit": "ns/op\t39689832 B/op\t  135382 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41879535,
            "unit": "ns/op\t39704816 B/op\t  135400 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41474997,
            "unit": "ns/op\t39730000 B/op\t  135435 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41563767,
            "unit": "ns/op\t39711682 B/op\t  135412 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38530094,
            "unit": "ns/op\t39659961 B/op\t  134967 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38939523,
            "unit": "ns/op\t39726616 B/op\t  135057 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38703066,
            "unit": "ns/op\t39644984 B/op\t  134946 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38707336,
            "unit": "ns/op\t39653549 B/op\t  134964 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38572713,
            "unit": "ns/op\t39699309 B/op\t  135018 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51677006,
            "unit": "ns/op\t42924362 B/op\t  139764 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52521891,
            "unit": "ns/op\t43070653 B/op\t  139863 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52006876,
            "unit": "ns/op\t42882794 B/op\t  139726 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52305836,
            "unit": "ns/op\t42950052 B/op\t  139771 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52316497,
            "unit": "ns/op\t42743392 B/op\t  139632 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43350874,
            "unit": "ns/op\t42220702 B/op\t  156368 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42595327,
            "unit": "ns/op\t42236321 B/op\t  156387 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42926590,
            "unit": "ns/op\t42243401 B/op\t  156407 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42889656,
            "unit": "ns/op\t42226635 B/op\t  156378 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42947772,
            "unit": "ns/op\t42237279 B/op\t  156385 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43373057,
            "unit": "ns/op\t42213834 B/op\t  156356 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43199754,
            "unit": "ns/op\t42202968 B/op\t  156338 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43077306,
            "unit": "ns/op\t42200331 B/op\t  156342 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43187670,
            "unit": "ns/op\t42216002 B/op\t  156370 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42996335,
            "unit": "ns/op\t42269184 B/op\t  156430 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 116001958,
            "unit": "ns/op\t62155819 B/op\t  715194 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115717100,
            "unit": "ns/op\t62549706 B/op\t  715351 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115188925,
            "unit": "ns/op\t62283915 B/op\t  715244 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115409249,
            "unit": "ns/op\t62007073 B/op\t  715135 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115945361,
            "unit": "ns/op\t62467781 B/op\t  715303 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benye@amazon.com",
            "name": "Ben Ye",
            "username": "yeya24"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "57db8dad11c2083cf8bf80617920305508aae3cd",
          "message": "add link to continuous benchmark (#165)\n\nSigned-off-by: Ben Ye <benye@amazon.com>",
          "timestamp": "2023-01-30T03:45:20-05:00",
          "tree_id": "8f1858fcab0a30e566b8b90598ca95616dd67814",
          "url": "https://github.com/thanos-community/promql-engine/commit/57db8dad11c2083cf8bf80617920305508aae3cd"
        },
        "date": 1675068489994,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22342621,
            "unit": "ns/op\t39362666 B/op\t  131326 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22202326,
            "unit": "ns/op\t39337435 B/op\t  131301 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23034993,
            "unit": "ns/op\t39397247 B/op\t  131363 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23022588,
            "unit": "ns/op\t39363228 B/op\t  131333 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23264240,
            "unit": "ns/op\t39378748 B/op\t  131344 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12155702,
            "unit": "ns/op\t11287201 B/op\t  125955 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12213372,
            "unit": "ns/op\t11277887 B/op\t  125965 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12233149,
            "unit": "ns/op\t11299825 B/op\t  125977 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12236420,
            "unit": "ns/op\t11299537 B/op\t  125974 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12270282,
            "unit": "ns/op\t11300727 B/op\t  125987 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17497636,
            "unit": "ns/op\t24565301 B/op\t  211773 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17488438,
            "unit": "ns/op\t24561726 B/op\t  211743 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17548352,
            "unit": "ns/op\t24565994 B/op\t  211754 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17594557,
            "unit": "ns/op\t24561224 B/op\t  211763 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17487567,
            "unit": "ns/op\t24561533 B/op\t  211744 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12779461,
            "unit": "ns/op\t12681380 B/op\t  132249 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12846231,
            "unit": "ns/op\t12731197 B/op\t  132295 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12783975,
            "unit": "ns/op\t12745059 B/op\t  132304 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12741762,
            "unit": "ns/op\t12726018 B/op\t  132275 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12853485,
            "unit": "ns/op\t12727825 B/op\t  132286 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12779641,
            "unit": "ns/op\t12768321 B/op\t  135984 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12583701,
            "unit": "ns/op\t12751662 B/op\t  135985 allocs/op",
            "extra": "87 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12921123,
            "unit": "ns/op\t12810132 B/op\t  136030 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12827639,
            "unit": "ns/op\t12798511 B/op\t  136017 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12768548,
            "unit": "ns/op\t12741002 B/op\t  135976 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28243512,
            "unit": "ns/op\t42194526 B/op\t  155221 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28416864,
            "unit": "ns/op\t42201024 B/op\t  155235 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28338202,
            "unit": "ns/op\t42211848 B/op\t  155248 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28533231,
            "unit": "ns/op\t42231666 B/op\t  155269 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28325634,
            "unit": "ns/op\t42188656 B/op\t  155223 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24143475,
            "unit": "ns/op\t13944517 B/op\t  149768 allocs/op",
            "extra": "73 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24020875,
            "unit": "ns/op\t14020724 B/op\t  149836 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24451461,
            "unit": "ns/op\t13960660 B/op\t  149773 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24456140,
            "unit": "ns/op\t13837151 B/op\t  149670 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24314622,
            "unit": "ns/op\t13890893 B/op\t  149709 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27690917,
            "unit": "ns/op\t26351596 B/op\t  235093 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27644674,
            "unit": "ns/op\t26331697 B/op\t  235078 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27628706,
            "unit": "ns/op\t26374946 B/op\t  235145 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27762138,
            "unit": "ns/op\t26406269 B/op\t  235145 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27727044,
            "unit": "ns/op\t26351128 B/op\t  235114 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26806924,
            "unit": "ns/op\t42899207 B/op\t  619155 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26780915,
            "unit": "ns/op\t42855798 B/op\t  619165 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26689469,
            "unit": "ns/op\t42861702 B/op\t  619161 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26787411,
            "unit": "ns/op\t42865520 B/op\t  619159 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26699369,
            "unit": "ns/op\t42876877 B/op\t  619168 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19477015,
            "unit": "ns/op\t23113476 B/op\t  108648 allocs/op",
            "extra": "64 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19394337,
            "unit": "ns/op\t23098748 B/op\t  108617 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19480122,
            "unit": "ns/op\t23099166 B/op\t  108621 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19581056,
            "unit": "ns/op\t23098625 B/op\t  108614 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19350670,
            "unit": "ns/op\t23103803 B/op\t  108625 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42380186,
            "unit": "ns/op\t58117374 B/op\t  202498 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42263823,
            "unit": "ns/op\t58129538 B/op\t  202554 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42286203,
            "unit": "ns/op\t58112074 B/op\t  202517 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42381126,
            "unit": "ns/op\t58133131 B/op\t  202554 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42118543,
            "unit": "ns/op\t58078105 B/op\t  202469 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33420896,
            "unit": "ns/op\t42692692 B/op\t  135585 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33215299,
            "unit": "ns/op\t42675096 B/op\t  135575 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33249260,
            "unit": "ns/op\t42605128 B/op\t  135529 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33277737,
            "unit": "ns/op\t42797433 B/op\t  135635 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33677931,
            "unit": "ns/op\t42650226 B/op\t  135565 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24563939,
            "unit": "ns/op\t40773197 B/op\t  143833 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24327723,
            "unit": "ns/op\t40780258 B/op\t  143827 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24320046,
            "unit": "ns/op\t40780179 B/op\t  143826 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24281850,
            "unit": "ns/op\t40789880 B/op\t  143842 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24493930,
            "unit": "ns/op\t40782616 B/op\t  143826 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33916461,
            "unit": "ns/op\t42390711 B/op\t  132612 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33518190,
            "unit": "ns/op\t42351047 B/op\t  132588 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33817608,
            "unit": "ns/op\t42313617 B/op\t  132567 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33404918,
            "unit": "ns/op\t42261657 B/op\t  132541 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33636548,
            "unit": "ns/op\t42279049 B/op\t  132543 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23040050,
            "unit": "ns/op\t37618094 B/op\t  102516 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22682961,
            "unit": "ns/op\t37610792 B/op\t  102515 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22806899,
            "unit": "ns/op\t37630554 B/op\t  102530 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22637495,
            "unit": "ns/op\t37622665 B/op\t  102523 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22496199,
            "unit": "ns/op\t37619425 B/op\t  102520 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17104810,
            "unit": "ns/op\t33220030 B/op\t   74763 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17391750,
            "unit": "ns/op\t33219382 B/op\t   74764 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17311379,
            "unit": "ns/op\t33220450 B/op\t   74763 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16854938,
            "unit": "ns/op\t33221207 B/op\t   74763 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16665208,
            "unit": "ns/op\t33221055 B/op\t   74763 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16923938,
            "unit": "ns/op\t33044448 B/op\t   68766 allocs/op",
            "extra": "80 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17250300,
            "unit": "ns/op\t33041041 B/op\t   68765 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17011774,
            "unit": "ns/op\t33040688 B/op\t   68765 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17221227,
            "unit": "ns/op\t33044779 B/op\t   68766 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17254193,
            "unit": "ns/op\t33045607 B/op\t   68766 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41179315,
            "unit": "ns/op\t39706682 B/op\t  135104 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 40740518,
            "unit": "ns/op\t39733441 B/op\t  135126 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41451795,
            "unit": "ns/op\t39670728 B/op\t  135046 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41575428,
            "unit": "ns/op\t39711980 B/op\t  135098 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41704984,
            "unit": "ns/op\t39710693 B/op\t  135101 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38626351,
            "unit": "ns/op\t39763967 B/op\t  134782 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38852923,
            "unit": "ns/op\t39633908 B/op\t  134615 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38811791,
            "unit": "ns/op\t39682688 B/op\t  134678 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38405806,
            "unit": "ns/op\t39766007 B/op\t  134789 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39110233,
            "unit": "ns/op\t39669514 B/op\t  134663 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52054939,
            "unit": "ns/op\t42896566 B/op\t  139431 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51511334,
            "unit": "ns/op\t42747697 B/op\t  139311 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52333706,
            "unit": "ns/op\t42745044 B/op\t  139334 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51817295,
            "unit": "ns/op\t42834704 B/op\t  139391 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51678589,
            "unit": "ns/op\t42968351 B/op\t  139472 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42331891,
            "unit": "ns/op\t42203069 B/op\t  156034 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42626726,
            "unit": "ns/op\t42201730 B/op\t  156021 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43002330,
            "unit": "ns/op\t42204403 B/op\t  156038 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42976204,
            "unit": "ns/op\t42249278 B/op\t  156084 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42784213,
            "unit": "ns/op\t42212620 B/op\t  156025 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42707672,
            "unit": "ns/op\t42212135 B/op\t  156040 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42924884,
            "unit": "ns/op\t42250293 B/op\t  156100 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42684270,
            "unit": "ns/op\t42223533 B/op\t  156052 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42596939,
            "unit": "ns/op\t42216140 B/op\t  156051 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42950170,
            "unit": "ns/op\t42220595 B/op\t  156053 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115147325,
            "unit": "ns/op\t62230747 B/op\t  715468 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113911833,
            "unit": "ns/op\t62106735 B/op\t  715435 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114030056,
            "unit": "ns/op\t62281696 B/op\t  715503 allocs/op",
            "extra": "10 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 116328607,
            "unit": "ns/op\t62262526 B/op\t  715511 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114060387,
            "unit": "ns/op\t61868645 B/op\t  715334 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benye@amazon.com",
            "name": "Ben Ye",
            "username": "yeya24"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e6c31f4efe0179026e932fe2c6ac20a453455fbe",
          "message": "handle no arg function labels (#162)\n\nSigned-off-by: Ben Ye <benye@amazon.com>",
          "timestamp": "2023-01-30T03:04:40-08:00",
          "tree_id": "67165ed3e456fe8da181f6f4fa1de71c262e57e3",
          "url": "https://github.com/thanos-community/promql-engine/commit/e6c31f4efe0179026e932fe2c6ac20a453455fbe"
        },
        "date": 1675076855302,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22443034,
            "unit": "ns/op\t39346920 B/op\t  131341 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22290897,
            "unit": "ns/op\t39345158 B/op\t  131339 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22478686,
            "unit": "ns/op\t39335601 B/op\t  131320 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22484350,
            "unit": "ns/op\t39343138 B/op\t  131330 allocs/op",
            "extra": "57 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23382409,
            "unit": "ns/op\t39350103 B/op\t  131339 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12403291,
            "unit": "ns/op\t11216085 B/op\t  125934 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12518156,
            "unit": "ns/op\t11237001 B/op\t  125930 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12524961,
            "unit": "ns/op\t11240899 B/op\t  125939 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12564661,
            "unit": "ns/op\t11290846 B/op\t  125989 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12552323,
            "unit": "ns/op\t11256857 B/op\t  125939 allocs/op",
            "extra": "94 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17785025,
            "unit": "ns/op\t24570850 B/op\t  211790 allocs/op",
            "extra": "72 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17751157,
            "unit": "ns/op\t24572782 B/op\t  211787 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17679856,
            "unit": "ns/op\t24576842 B/op\t  211800 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17733979,
            "unit": "ns/op\t24575105 B/op\t  211793 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17700401,
            "unit": "ns/op\t24585000 B/op\t  211803 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12984099,
            "unit": "ns/op\t12673965 B/op\t  133445 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12960425,
            "unit": "ns/op\t12612559 B/op\t  133379 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13030681,
            "unit": "ns/op\t12637522 B/op\t  133404 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12960105,
            "unit": "ns/op\t12651541 B/op\t  133422 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12952205,
            "unit": "ns/op\t12608696 B/op\t  133374 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12984459,
            "unit": "ns/op\t12676178 B/op\t  134672 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12983570,
            "unit": "ns/op\t12688653 B/op\t  134672 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 13012162,
            "unit": "ns/op\t12638114 B/op\t  134631 allocs/op",
            "extra": "90 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 13013499,
            "unit": "ns/op\t12646022 B/op\t  134648 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 13044079,
            "unit": "ns/op\t12682348 B/op\t  134672 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28668251,
            "unit": "ns/op\t42215953 B/op\t  155287 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28728110,
            "unit": "ns/op\t42188451 B/op\t  155245 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28773812,
            "unit": "ns/op\t42196185 B/op\t  155266 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28663706,
            "unit": "ns/op\t42137945 B/op\t  155195 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28879048,
            "unit": "ns/op\t42155855 B/op\t  155211 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24673995,
            "unit": "ns/op\t13857482 B/op\t  149712 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24712051,
            "unit": "ns/op\t13904712 B/op\t  149763 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24549356,
            "unit": "ns/op\t13993341 B/op\t  149807 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24631211,
            "unit": "ns/op\t13829850 B/op\t  149689 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24563682,
            "unit": "ns/op\t13964817 B/op\t  149805 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28122341,
            "unit": "ns/op\t26364209 B/op\t  235098 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28130069,
            "unit": "ns/op\t26303372 B/op\t  235087 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28132352,
            "unit": "ns/op\t26396574 B/op\t  235147 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 28099557,
            "unit": "ns/op\t26310222 B/op\t  235094 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27997022,
            "unit": "ns/op\t26371627 B/op\t  235130 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27220569,
            "unit": "ns/op\t42839354 B/op\t  619185 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27035532,
            "unit": "ns/op\t42865794 B/op\t  619218 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27035795,
            "unit": "ns/op\t42860963 B/op\t  619187 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 27022248,
            "unit": "ns/op\t42845306 B/op\t  619202 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26955396,
            "unit": "ns/op\t42849581 B/op\t  619223 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19604669,
            "unit": "ns/op\t23097781 B/op\t  108543 allocs/op",
            "extra": "68 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19653906,
            "unit": "ns/op\t23091573 B/op\t  108500 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19795945,
            "unit": "ns/op\t23127539 B/op\t  108590 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19527139,
            "unit": "ns/op\t23093803 B/op\t  108508 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19594846,
            "unit": "ns/op\t23091327 B/op\t  108524 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43166553,
            "unit": "ns/op\t58087912 B/op\t  202524 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43103596,
            "unit": "ns/op\t58167198 B/op\t  202638 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43109916,
            "unit": "ns/op\t58111421 B/op\t  202547 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43359771,
            "unit": "ns/op\t58074147 B/op\t  202519 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 43391885,
            "unit": "ns/op\t58115901 B/op\t  202579 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33752232,
            "unit": "ns/op\t42678064 B/op\t  135599 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33298524,
            "unit": "ns/op\t42688078 B/op\t  135599 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33615814,
            "unit": "ns/op\t42598687 B/op\t  135547 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 34020311,
            "unit": "ns/op\t42631824 B/op\t  135566 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33262298,
            "unit": "ns/op\t42689308 B/op\t  135602 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24630910,
            "unit": "ns/op\t40726445 B/op\t  143783 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24786921,
            "unit": "ns/op\t40733198 B/op\t  143794 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24545639,
            "unit": "ns/op\t40723040 B/op\t  143800 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24703041,
            "unit": "ns/op\t40743923 B/op\t  143798 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24439909,
            "unit": "ns/op\t40728059 B/op\t  143782 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33867210,
            "unit": "ns/op\t42133876 B/op\t  132485 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33838483,
            "unit": "ns/op\t42343474 B/op\t  132605 allocs/op",
            "extra": "36 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33291009,
            "unit": "ns/op\t42293191 B/op\t  132574 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33627474,
            "unit": "ns/op\t42221739 B/op\t  132544 allocs/op",
            "extra": "39 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34283384,
            "unit": "ns/op\t42216441 B/op\t  132528 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22971603,
            "unit": "ns/op\t37604033 B/op\t  102529 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23093867,
            "unit": "ns/op\t37595905 B/op\t  102509 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22989554,
            "unit": "ns/op\t37563527 B/op\t  102469 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22982311,
            "unit": "ns/op\t37597354 B/op\t  102514 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 23098453,
            "unit": "ns/op\t37590740 B/op\t  102508 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17072898,
            "unit": "ns/op\t33214340 B/op\t   74769 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17532004,
            "unit": "ns/op\t33217871 B/op\t   74770 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17215585,
            "unit": "ns/op\t33217881 B/op\t   74769 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17280943,
            "unit": "ns/op\t33217104 B/op\t   74769 allocs/op",
            "extra": "79 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17393750,
            "unit": "ns/op\t33216941 B/op\t   74769 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16909066,
            "unit": "ns/op\t33035432 B/op\t   68771 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17395275,
            "unit": "ns/op\t33038582 B/op\t   68771 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16910894,
            "unit": "ns/op\t33041335 B/op\t   68772 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16860832,
            "unit": "ns/op\t33038392 B/op\t   68771 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16939878,
            "unit": "ns/op\t33039759 B/op\t   68771 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41336530,
            "unit": "ns/op\t39654034 B/op\t  135055 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41715821,
            "unit": "ns/op\t39661632 B/op\t  135059 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41084599,
            "unit": "ns/op\t39652430 B/op\t  135039 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42199039,
            "unit": "ns/op\t39667162 B/op\t  135070 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 42068784,
            "unit": "ns/op\t39646748 B/op\t  135042 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39059952,
            "unit": "ns/op\t39612250 B/op\t  134624 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38888885,
            "unit": "ns/op\t39655097 B/op\t  134678 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38632777,
            "unit": "ns/op\t39651327 B/op\t  134668 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38644411,
            "unit": "ns/op\t39677597 B/op\t  134707 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39108459,
            "unit": "ns/op\t39674777 B/op\t  134705 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52432565,
            "unit": "ns/op\t42828361 B/op\t  139394 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52246400,
            "unit": "ns/op\t42966336 B/op\t  139487 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52630617,
            "unit": "ns/op\t43018657 B/op\t  139520 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52795917,
            "unit": "ns/op\t42809751 B/op\t  139384 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52573853,
            "unit": "ns/op\t42851173 B/op\t  139410 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42911366,
            "unit": "ns/op\t42225877 B/op\t  156088 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43436366,
            "unit": "ns/op\t42206650 B/op\t  156059 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43480191,
            "unit": "ns/op\t42201229 B/op\t  156063 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43110370,
            "unit": "ns/op\t42200974 B/op\t  156051 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43071193,
            "unit": "ns/op\t42269936 B/op\t  156161 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43783562,
            "unit": "ns/op\t42181349 B/op\t  156021 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42355094,
            "unit": "ns/op\t42227136 B/op\t  156097 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43483576,
            "unit": "ns/op\t42212954 B/op\t  156066 allocs/op",
            "extra": "27 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43230177,
            "unit": "ns/op\t42216138 B/op\t  156070 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43038145,
            "unit": "ns/op\t42218035 B/op\t  156080 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114159648,
            "unit": "ns/op\t62235342 B/op\t  715310 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115812478,
            "unit": "ns/op\t62473383 B/op\t  715414 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114746719,
            "unit": "ns/op\t62314586 B/op\t  715337 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115490416,
            "unit": "ns/op\t62392922 B/op\t  715375 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115063715,
            "unit": "ns/op\t62526828 B/op\t  715421 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "121007071+kama910@users.noreply.github.com",
            "name": "Kama Huang",
            "username": "kama910"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "65156b30319d030aa1184722becf967e502d00ca",
          "message": "add support for a few functions (#154)\n\n* add support for a few functions\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>\r\n\r\n* support variadic functions\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>\r\n\r\n---------\r\n\r\nSigned-off-by: Kama Huang <kamatogo13@gmail.com>",
          "timestamp": "2023-01-31T10:04:49+02:00",
          "tree_id": "ee0fa2bafd1c15225e972e9da041a24f6af3cf83",
          "url": "https://github.com/thanos-community/promql-engine/commit/65156b30319d030aa1184722becf967e502d00ca"
        },
        "date": 1675152462509,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22379984,
            "unit": "ns/op\t39352124 B/op\t  131730 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 22663698,
            "unit": "ns/op\t39375279 B/op\t  131752 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23336159,
            "unit": "ns/op\t39359921 B/op\t  131737 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23351748,
            "unit": "ns/op\t39381128 B/op\t  131764 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_selector",
            "value": 23287684,
            "unit": "ns/op\t39396644 B/op\t  131778 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12472987,
            "unit": "ns/op\t11311762 B/op\t  126403 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12462565,
            "unit": "ns/op\t11258771 B/op\t  126337 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12476077,
            "unit": "ns/op\t11277100 B/op\t  126354 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12513042,
            "unit": "ns/op\t11287006 B/op\t  126371 allocs/op",
            "extra": "92 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum",
            "value": 12515929,
            "unit": "ns/op\t11267870 B/op\t  126361 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17569328,
            "unit": "ns/op\t24556068 B/op\t  212154 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17744858,
            "unit": "ns/op\t24583924 B/op\t  212186 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17594005,
            "unit": "ns/op\t24593064 B/op\t  212193 allocs/op",
            "extra": "70 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17695408,
            "unit": "ns/op\t24608400 B/op\t  212206 allocs/op",
            "extra": "69 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_pod",
            "value": 17656562,
            "unit": "ns/op\t24589943 B/op\t  212192 allocs/op",
            "extra": "67 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13016139,
            "unit": "ns/op\t12762248 B/op\t  135503 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13050421,
            "unit": "ns/op\t12739033 B/op\t  135469 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13000657,
            "unit": "ns/op\t12700839 B/op\t  135441 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 12982263,
            "unit": "ns/op\t12772472 B/op\t  135497 allocs/op",
            "extra": "93 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/topk",
            "value": 13005797,
            "unit": "ns/op\t12667171 B/op\t  135421 allocs/op",
            "extra": "87 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12938143,
            "unit": "ns/op\t12681525 B/op\t  133345 allocs/op",
            "extra": "96 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 13004522,
            "unit": "ns/op\t12651066 B/op\t  133329 allocs/op",
            "extra": "85 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12995189,
            "unit": "ns/op\t12714889 B/op\t  133383 allocs/op",
            "extra": "88 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12975808,
            "unit": "ns/op\t12674714 B/op\t  133344 allocs/op",
            "extra": "91 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/bottomk",
            "value": 12963963,
            "unit": "ns/op\t12697312 B/op\t  133356 allocs/op",
            "extra": "87 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28500127,
            "unit": "ns/op\t42235684 B/op\t  155689 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28665624,
            "unit": "ns/op\t42206282 B/op\t  155657 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28470189,
            "unit": "ns/op\t42209684 B/op\t  155659 allocs/op",
            "extra": "42 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28424768,
            "unit": "ns/op\t42178748 B/op\t  155631 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/rate",
            "value": 28517644,
            "unit": "ns/op\t42241331 B/op\t  155693 allocs/op",
            "extra": "44 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24499258,
            "unit": "ns/op\t13898664 B/op\t  150143 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24458095,
            "unit": "ns/op\t13916280 B/op\t  150167 allocs/op",
            "extra": "49 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24526011,
            "unit": "ns/op\t13928459 B/op\t  150159 allocs/op",
            "extra": "45 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24532122,
            "unit": "ns/op\t13895270 B/op\t  150129 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_rate",
            "value": 24612920,
            "unit": "ns/op\t13864808 B/op\t  150107 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27898531,
            "unit": "ns/op\t26367463 B/op\t  235514 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27908114,
            "unit": "ns/op\t26381202 B/op\t  235532 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27912462,
            "unit": "ns/op\t26421367 B/op\t  235588 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27903754,
            "unit": "ns/op\t26328686 B/op\t  235500 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/sum_by_rate",
            "value": 27716827,
            "unit": "ns/op\t26356382 B/op\t  235507 allocs/op",
            "extra": "43 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26987174,
            "unit": "ns/op\t42912020 B/op\t  619998 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26884290,
            "unit": "ns/op\t42921842 B/op\t  620080 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26872167,
            "unit": "ns/op\t42904628 B/op\t  620045 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26733198,
            "unit": "ns/op\t42893486 B/op\t  620023 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/quantile_with_variable_parameter",
            "value": 26847479,
            "unit": "ns/op\t42895557 B/op\t  620015 allocs/op",
            "extra": "46 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19554264,
            "unit": "ns/op\t23124016 B/op\t  108565 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19424519,
            "unit": "ns/op\t23103588 B/op\t  108526 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19636040,
            "unit": "ns/op\t23101519 B/op\t  108530 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19435996,
            "unit": "ns/op\t23132263 B/op\t  108606 allocs/op",
            "extra": "61 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_one_to_one",
            "value": 19640791,
            "unit": "ns/op\t23112624 B/op\t  108572 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42960544,
            "unit": "ns/op\t58083584 B/op\t  202561 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42812829,
            "unit": "ns/op\t58133952 B/op\t  202641 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42526175,
            "unit": "ns/op\t58133309 B/op\t  202644 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42333841,
            "unit": "ns/op\t58146128 B/op\t  202676 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_many_to_one",
            "value": 42582769,
            "unit": "ns/op\t58083063 B/op\t  202570 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33548287,
            "unit": "ns/op\t42825101 B/op\t  136083 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33795789,
            "unit": "ns/op\t42725271 B/op\t  136020 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33026372,
            "unit": "ns/op\t42643818 B/op\t  135978 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33555428,
            "unit": "ns/op\t42698761 B/op\t  136016 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/binary_operation_with_vector_and_scalar",
            "value": 33839436,
            "unit": "ns/op\t42558745 B/op\t  135924 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24384381,
            "unit": "ns/op\t40760263 B/op\t  144213 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24733705,
            "unit": "ns/op\t40756576 B/op\t  144206 allocs/op",
            "extra": "51 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24838443,
            "unit": "ns/op\t40775653 B/op\t  144230 allocs/op",
            "extra": "48 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24593011,
            "unit": "ns/op\t40761266 B/op\t  144201 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/unary_negation",
            "value": 24692648,
            "unit": "ns/op\t40764294 B/op\t  144199 allocs/op",
            "extra": "50 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34023485,
            "unit": "ns/op\t42221280 B/op\t  132944 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 34068648,
            "unit": "ns/op\t42303249 B/op\t  132979 allocs/op",
            "extra": "34 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33752212,
            "unit": "ns/op\t42219154 B/op\t  132938 allocs/op",
            "extra": "38 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33732795,
            "unit": "ns/op\t42167838 B/op\t  132904 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/vector_and_scalar_comparison",
            "value": 33691941,
            "unit": "ns/op\t42287096 B/op\t  132976 allocs/op",
            "extra": "37 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22832617,
            "unit": "ns/op\t37593379 B/op\t  102765 allocs/op",
            "extra": "52 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22678984,
            "unit": "ns/op\t37605577 B/op\t  102775 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22668706,
            "unit": "ns/op\t37645870 B/op\t  102817 allocs/op",
            "extra": "56 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22782308,
            "unit": "ns/op\t37625137 B/op\t  102794 allocs/op",
            "extra": "54 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/positive_offset_vector",
            "value": 22598988,
            "unit": "ns/op\t37630532 B/op\t  102805 allocs/op",
            "extra": "55 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17456240,
            "unit": "ns/op\t33223286 B/op\t   74906 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17033900,
            "unit": "ns/op\t33224848 B/op\t   74906 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 16910698,
            "unit": "ns/op\t33225403 B/op\t   74907 allocs/op",
            "extra": "81 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17025422,
            "unit": "ns/op\t33224220 B/op\t   74907 allocs/op",
            "extra": "62 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_",
            "value": 17451764,
            "unit": "ns/op\t33223796 B/op\t   74907 allocs/op",
            "extra": "63 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16687318,
            "unit": "ns/op\t33047616 B/op\t   68910 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17200340,
            "unit": "ns/op\t33050708 B/op\t   68910 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16641506,
            "unit": "ns/op\t33049957 B/op\t   68909 allocs/op",
            "extra": "84 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 17421449,
            "unit": "ns/op\t33048443 B/op\t   68911 allocs/op",
            "extra": "66 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/at_modifier_with_positive_offset_vector",
            "value": 16960460,
            "unit": "ns/op\t33050178 B/op\t   68910 allocs/op",
            "extra": "82 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41804503,
            "unit": "ns/op\t39696932 B/op\t  135482 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41788472,
            "unit": "ns/op\t39647703 B/op\t  135413 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41667992,
            "unit": "ns/op\t39668360 B/op\t  135446 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41932921,
            "unit": "ns/op\t39656350 B/op\t  135421 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp",
            "value": 41739543,
            "unit": "ns/op\t39692446 B/op\t  135488 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39069021,
            "unit": "ns/op\t39673344 B/op\t  135076 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38936487,
            "unit": "ns/op\t39654675 B/op\t  135056 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 38724778,
            "unit": "ns/op\t39686699 B/op\t  135093 allocs/op",
            "extra": "32 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39236268,
            "unit": "ns/op\t39640307 B/op\t  135037 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/clamp_min",
            "value": 39040560,
            "unit": "ns/op\t39627706 B/op\t  135016 allocs/op",
            "extra": "31 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52237323,
            "unit": "ns/op\t42803736 B/op\t  139769 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 52465481,
            "unit": "ns/op\t42909839 B/op\t  139856 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51641979,
            "unit": "ns/op\t42800555 B/op\t  139774 allocs/op",
            "extra": "22 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 50953612,
            "unit": "ns/op\t42854203 B/op\t  139817 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/complex_func_query",
            "value": 51988311,
            "unit": "ns/op\t42850503 B/op\t  139808 allocs/op",
            "extra": "24 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 43194026,
            "unit": "ns/op\t42205915 B/op\t  156432 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42922255,
            "unit": "ns/op\t42235566 B/op\t  156477 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42849780,
            "unit": "ns/op\t42280980 B/op\t  156541 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42868792,
            "unit": "ns/op\t42208600 B/op\t  156432 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/func_within_func_query",
            "value": 42949213,
            "unit": "ns/op\t42208976 B/op\t  156441 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43426293,
            "unit": "ns/op\t42236774 B/op\t  156476 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43118291,
            "unit": "ns/op\t42182735 B/op\t  156411 allocs/op",
            "extra": "30 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42973828,
            "unit": "ns/op\t42214304 B/op\t  156439 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 43149383,
            "unit": "ns/op\t42169266 B/op\t  156374 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/aggr_within_func_query",
            "value": 42895925,
            "unit": "ns/op\t42212879 B/op\t  156453 allocs/op",
            "extra": "28 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115164279,
            "unit": "ns/op\t62372123 B/op\t  715387 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 113812490,
            "unit": "ns/op\t61915457 B/op\t  715207 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114561381,
            "unit": "ns/op\t62056670 B/op\t  715277 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 114967672,
            "unit": "ns/op\t62358331 B/op\t  715382 allocs/op",
            "extra": "9 times\n16 procs"
          },
          {
            "name": "BenchmarkRangeQuery/histogram_quantile",
            "value": 115787621,
            "unit": "ns/op\t62101352 B/op\t  715296 allocs/op",
            "extra": "9 times\n16 procs"
          }
        ]
      }
    ]
  }
}