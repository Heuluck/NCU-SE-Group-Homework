import React from "react";
import { Pagination as AntdPagination } from "antd";

type Props = {
  current: number;
  pageSize: number;
  total: number;
  onChange: (page: number, pageSize?: number) => void;
};

export default function Pagination({ current, pageSize, total, onChange }: Props) {
  return (
    <div style={{ display: "flex", justifyContent: "flex-end", paddingTop: 8 }}>
      <AntdPagination
        current={current}
        pageSize={pageSize}
        total={total}
        onChange={onChange}
        showSizeChanger
        pageSizeOptions={["5", "10", "20"]}
      />
    </div>
  );
}
