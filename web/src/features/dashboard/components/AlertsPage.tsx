import { Badge, Button, DataTable, Panel, type DataColumn } from "@netstamp/ui";
import { ScreenHeader } from "../../../shared/components/ScreenHeader";
import { alerts, toneForStatus, type AlertRecord } from "../../../shared/utils/mockData";
import styles from "./ProductPages.module.css";

const alertColumns: DataColumn<AlertRecord>[] = [
	{ key: "type", label: "Alert type" },
	{ key: "check", label: "Check" },
	{ key: "probe", label: "Probe" },
	{ key: "severity", label: "Severity", render: row => <Badge tone={toneForStatus(row.severity)}>{row.severity}</Badge> },
	{ key: "triggered", label: "Triggered" },
	{ key: "status", label: "Status", render: row => <Badge tone={toneForStatus(row.status)}>{row.status}</Badge> }
];

export function AlertsPage() {
	return (
		<section className={styles.screen}>
			<ScreenHeader eyebrow="Alerting" title="Alerts (TBD)" copy="Packet loss, latency, traceroute path change, DNS query errors, abnormal response codes, probe offline, and heartbeat expiry." />

			<div className={styles.alertGrid}>
				<Panel tone="glass" eyebrow="Alert list" title="Active and historical events">
					<DataTable columns={alertColumns} rows={alerts} getRowKey={row => `${row.type}-${row.probe}`} />
				</Panel>
				<Panel tone="deep" eyebrow="Alert detail" title="packet loss threshold exceeded">
					<div className={styles.keyValueGrid}>
						<div>
							<span>Affected probe</span>
							<strong>nyc-vps-03</strong>
						</div>
						<div>
							<span>Affected check</span>
							<strong>api-latency</strong>
						</div>
						<div>
							<span>Threshold</span>
							<strong>loss &gt; 5% for 5m</strong>
						</div>
						<div>
							<span>State</span>
							<strong>active</strong>
						</div>
					</div>
					<div className={styles.actionRow}>
						<Button variant="secondary">Open result history</Button>
						<Button variant="danger">Silence 30m</Button>
					</div>
				</Panel>
			</div>
		</section>
	);
}
