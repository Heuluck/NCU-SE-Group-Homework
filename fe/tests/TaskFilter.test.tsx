import React from "react";
import { render, screen } from "@testing-library/react";
import TaskFilter from "../app/components/TaskFilter";

describe("TaskFilter", () => {
  it("renders options", () => {
    render(<TaskFilter />);
    expect(screen.getByText("全部")).toBeInTheDocument();
    expect(screen.getByText("待办")).toBeInTheDocument();
    expect(screen.getByText("已完成")).toBeInTheDocument();
  });
});
