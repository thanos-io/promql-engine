window.BENCHMARK_DATA = {
  "lastUpdate": 1668540443674,
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
      }
    ]
  }
}