import React from "react";
import { Select } from "antd";

const { Option } = Select;

type Props = {
  value?: string;
  onChange?: (val: string) => void;
};

export default function TaskFilter({ value = "all", onChange }: Props) {
  return (
    <Select value={value} onChange={onChange} style={{ width: 160 }}>
      <Option value="all">全部</Option>
      <Option value="pending">待办</Option>
      <Option value="completed">已完成</Option>
    </Select>
  );
}
