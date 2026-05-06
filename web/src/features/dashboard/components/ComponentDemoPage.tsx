import {
	Badge,
	Button,
	DataTable,
	MetricCard,
	Panel,
	SelectField,
	Terminal,
	TextAreaField,
	TextField,
	type BadgeTone,
	type ButtonSize,
	type ButtonVariant,
	type DataColumn,
	type PanelTone
} from "@netstamp/ui";
import { useState } from "react";
import { Link } from "react-router-dom";
import { pathForRoute } from "../../../routes/routePaths";
import { ChartPanel } from "../../../shared/components/ChartPanel";
import { FleetMatrix } from "../../../shared/components/FleetMatrix";
import { NetworkMap } from "../../../shared/components/NetworkMap";
import { ScreenHeader } from "../../../shared/components/ScreenHeader";
import { SystemStateGrid } from "../../../shared/components/SystemStateGrid";
import { barChartOption, lineChartOption } from "../../../shared/utils/chartOptions";
import { latencyData, lossData, probes, toneForStatus } from "../../../shared/utils/mockData";
import styles from "./ComponentDemoPage.module.css";

interface ComponentRow {
	id: string;
	name: string;
	layer: string;
	status: string;
	notes: string;
}

const badgeTones: BadgeTone[] = ["neutral", "accent", "success", "warning", "critical", "muted"];
const buttonVariants: ButtonVariant[] = ["primary", "secondary", "outline", "ghost", "danger"];
const buttonSizes: ButtonSize[] = ["sm", "md", "lg", "xl"];
const panelTones: PanelTone[] = ["glass", "matte", "deep"];

const componentRows: ComponentRow[] = [
	{ id: "button", name: "Button", layer: "@netstamp/ui", status: "Healthy", notes: "variants, sizes, link slot" },
	{ id: "badge", name: "Badge", layer: "@netstamp/ui", status: "Healthy", notes: "six tones with optional dot" },
	{ id: "panel", name: "Panel", layer: "@netstamp/ui", status: "Selected", notes: "glass, matte, deep surfaces" },
	{ id: "field", name: "TextField / SelectField / TextAreaField", layer: "@netstamp/ui", status: "Warning", notes: "helper and validation states" },
	{ id: "table", name: "DataTable", layer: "@netstamp/ui", status: "Healthy", notes: "selected row and custom renderers" },
	{ id: "telemetry", name: "ChartPanel / FleetMatrix / NetworkMap", layer: "app shared", status: "Healthy", notes: "operational data visuals" }
];

const componentColumns: DataColumn<ComponentRow>[] = [
	{ key: "name", label: "Component" },
	{ key: "layer", label: "Layer" },
	{ key: "status", label: "State", render: row => <Badge tone={toneForStatus(row.status)}>{row.status}</Badge> },
	{ key: "notes", label: "Coverage" }
];

const teamOptions = [
	{ value: "vector-ix", label: "Vector IX / prod" },
	{ value: "helio", label: "Helio Validators" },
	{ value: "lab", label: "Lab Network" }
];

export function ComponentDemoPage() {
	const [selectedRow, setSelectedRow] = useState("panel");
	const [selectedProbeId, setSelectedProbeId] = useState(probes[0].id);

	return (
		<section className={styles.screen}>
			<ScreenHeader
				eyebrow="Component system"
				title="Netstamp Component Demo"
				copy="One route for checking every reusable UI primitive and the app-level telemetry widgets against the production console theme."
				actions={
					<>
						<Button variant="secondary" asChild>
							<Link to={pathForRoute("dashboard")}>Back to console</Link>
						</Button>
						<Button asChild>
							<Link to={pathForRoute("settings")}>Open settings</Link>
						</Button>
					</>
				}
			/>

			<div className={styles.metricGrid}>
				<MetricCard label="UI exports" value="7" detail="package" tone="accent" />
				<MetricCard label="App widgets" value="5" detail="shared" tone="success" />
				<MetricCard label="Route mode" value="React Router" detail="hash" tone="warning" />
			</div>

			<div className={styles.demoGrid}>
				<Panel tone="glass" eyebrow="Buttons" title="Variants and sizes">
					<div className={styles.buttonMatrix}>
						{buttonVariants.map(variant => (
							<Button key={variant} type="button" variant={variant}>
								{variant}
							</Button>
						))}
					</div>
					<div className={styles.buttonMatrix}>
						{buttonSizes.map(size => (
							<Button key={size} type="button" size={size} variant="secondary">
								{size}
							</Button>
						))}
						<Button asChild variant="outline">
							<Link to={pathForRoute("probes")}>asChild link</Link>
						</Button>
					</div>
				</Panel>

				<Panel tone="glass" eyebrow="Badges" title="Tone scale">
					<div className={styles.badgeCluster}>
						{badgeTones.map(tone => (
							<Badge key={tone} tone={tone}>
								{tone}
							</Badge>
						))}
						<Badge tone="accent" dot={false}>
							no dot
						</Badge>
					</div>
				</Panel>

				<Panel className={styles.wide} tone="matte" eyebrow="Panels" title="Surface tones">
					<div className={styles.panelToneGrid}>
						{panelTones.map(tone => (
							<Panel key={tone} tone={tone} eyebrow={`${tone} tone`} title={`${tone} panel`}>
								<p className={styles.copy}>Reusable framed section with title, eyebrow, separator, actions, and tone-specific depth.</p>
							</Panel>
						))}
					</div>
				</Panel>

				<Panel tone="glass" eyebrow="Fields" title="Inputs and validation">
					<div className={styles.formGrid}>
						<TextField label="Probe name" defaultValue="ams-edge-01" helper="Text input with helper copy." />
						<SelectField label="Team" defaultValue="vector-ix" options={teamOptions} />
						<TextField label="Target" defaultValue="api.netstamp.io" error="Target is already used by another check." />
						<TextAreaField label="Runbook note" defaultValue="Investigate route changes before acknowledging the active packet loss alert." helper="Textarea uses the same field shell." />
					</div>
				</Panel>

				<Panel tone="deep" eyebrow="Terminal" title="Command surface">
					<Terminal title="controller shell" meta="demo">
						{`netstamp checks list --status active
netstamp probes watch --team vector-ix
netstamp results tail --check api-latency`}
					</Terminal>
				</Panel>

				<Panel className={styles.wide} tone="glass" eyebrow="Data table" title="Rows, selection, custom renderers">
					<DataTable columns={componentColumns} rows={componentRows} getRowKey={row => row.id} selectedKey={selectedRow} onRowClick={row => setSelectedRow(row.id)} />
				</Panel>

				<Panel tone="glass" eyebrow="Charts" title="Latency and packet loss">
					<div className={styles.chartGrid}>
						<ChartPanel option={lineChartOption("latency", latencyData)} height="12rem" />
						<ChartPanel
							option={barChartOption(
								lossData.map(value => value * 100),
								"loss"
							)}
							height="12rem"
						/>
					</div>
				</Panel>

				<Panel tone="deep" eyebrow="Fleet" title="Map and bitmap widgets">
					<div className={styles.widgetStack}>
						<FleetMatrix total={64} online={51} />
						<NetworkMap probes={probes} selectedId={selectedProbeId} onSelect={setSelectedProbeId} />
					</div>
				</Panel>

				<Panel className={styles.wide} tone="glass" eyebrow="State grid" title="Reusable loading, empty, permission, and error states">
					<SystemStateGrid />
				</Panel>
			</div>
		</section>
	);
}
