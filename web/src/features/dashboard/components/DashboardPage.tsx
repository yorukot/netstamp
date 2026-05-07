import { Badge, Button, MetricCard, Panel, type BadgeTone } from "@netstamp/ui";
import { FleetMatrix } from "../../../shared/components/FleetMatrix";
import { PageStack } from "../../../shared/components/PageStack";
import { ScreenHeader } from "../../../shared/components/ScreenHeader";
import { type Navigate } from "../../../shared/utils/mockData";
import styles from "./DashboardPage.module.css";

interface DashboardPageProps {
	navigate: Navigate;
}

export function DashboardPage({ navigate }: DashboardPageProps) {
	return (
		<PageStack>
			<ScreenHeader
				eyebrow="Controller overview"
				title="Dashboard"
				copy="A premium active-measurement cockpit for probe health, scheduled checks, stream semantics, and recent path anomalies."
				actions={
					<>
						<Button variant="secondary" onClick={() => navigate("probes", "#new-probe")}>
							New probe wizard
						</Button>
						<Button onClick={() => navigate("checks")}>Create check</Button>
					</>
				}
			/>

			<div className={styles.metricsGrid}>
				<MetricCard label="Probes Online" value="100/128" detail="fleet" tone="success" />
				<MetricCard label="Active Checks" value="324" detail="scheduled" tone="accent" />
			</div>

			<div className={styles.dashboardGrid}>
				<Panel tone="glass" eyebrow="Fleet bitmap" title="128 probes, 100 lit">
					<FleetMatrix total={128} online={100} />
				</Panel>
				<Panel tone="glass" eyebrow="Anomalies" title="Recent system events">
					<div className={styles.feed}>
						<Event title="Packet loss threshold exceeded" copy="nyc-vps-03 → api.netstamp.io exceeded 18% loss for 5m." tone="critical" />
						<Event title="Path hash changed from previous run" copy="fra-bm-02 observed transit shift at hop 9." tone="warning" />
						<Event title="Controller stream connected" copy="100 probes streaming normalized result payloads." tone="success" />
					</div>
				</Panel>
			</div>
		</PageStack>
	);
}

interface EventProps {
	title: string;
	copy: string;
	tone: BadgeTone;
}

function Event({ title, copy, tone }: EventProps) {
	return (
		<article className={styles.event}>
			<Badge tone={tone}>{title}</Badge>
			<p>{copy}</p>
		</article>
	);
}
