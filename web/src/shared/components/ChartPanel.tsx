import { BarChart, LineChart } from "echarts/charts";
import { GridComponent, LegendComponent, TooltipComponent } from "echarts/components";
import * as echarts from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { useEffect, useRef } from "react";
import type { ChartOption } from "../utils/chartOptions";
import styles from "./ChartPanel.module.css";

echarts.use([LineChart, BarChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer]);

interface ChartPanelProps {
	option: ChartOption;
	height?: string;
	className?: string;
}

type SetOptionValue = Parameters<ReturnType<typeof echarts.init>["setOption"]>[0];

export function ChartPanel({ option, height = "16rem", className }: ChartPanelProps) {
	const chartRef = useRef<HTMLDivElement | null>(null);

	useEffect(() => {
		if (!chartRef.current) {
			return undefined;
		}

		const chart = echarts.init(chartRef.current, null, { renderer: "canvas" });
		chart.setOption(option as SetOptionValue);

		const observer = new ResizeObserver(() => chart.resize());
		observer.observe(chartRef.current);

		return () => {
			observer.disconnect();
			chart.dispose();
		};
	}, [option]);

	return <div ref={chartRef} className={[styles.chart, className].filter(Boolean).join(" ")} style={{ height }} />;
}
