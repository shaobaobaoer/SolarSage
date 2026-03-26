#!/bin/bash
# 分析时间差异模式

echo "=== 时间差异模式分析 ==="
echo ""
echo "SF事件时间 vs 计算事件时间对比:"
echo ""

# 提取前10个未匹配的SF事件
echo "前10个未匹配的SF事件:"
grep "DEBUG: Exact key missing" ultraprecise_output.txt | head -10 | while read line; do
    # 提取日期时间
    datetime=$(echo "$line" | grep -oP '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}')
    echo "  SF: $datetime"
done

echo ""
echo "前10个计算事件:"
grep "DEBUG COMP KEY:" ultraprecise_output.txt | head -10 | while read line; do
    # 提取日期时间
    datetime=$(echo "$line" | grep -oP '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}')
    echo "  Computed: $datetime"
done