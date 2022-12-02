window.BENCHMARK_DATA = {
  "lastUpdate": 1669979301023,
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
      }
    ]
  }
}