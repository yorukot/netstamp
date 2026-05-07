import { TextField } from "@netstamp/ui";
import { useState } from "react";
import { classNames } from "../../../shared/utils/classNames";
import type { Probe } from "../../../shared/utils/mockData";
import styles from "./ProbeDetail.module.css";
import { expandAssignedRows } from "./probeUtils";
import type { AssignedRow, DetectionMode } from "./types";

interface ProbeDetailProps {
	probe: Probe;
	assignedRows: AssignedRow[];
	floating?: boolean;
}

export function ProbeDetail({ probe, assignedRows, floating = false }: ProbeDetailProps) {
	const [probeName, setProbeName] = useState(probe.name);
	const [probeLocation, setProbeLocation] = useState(probe.location);
	const [probeAsn, setProbeAsn] = useState(probe.asn);
	const [locationMode, setLocationMode] = useState<DetectionMode>("manual");
	const [asMode, setAsMode] = useState<DetectionMode>("auto");
	const probeAssignments = assignedRows.filter(row => row.probe === probe.name);
	const baseRows = probeAssignments.length ? probeAssignments : assignedRows.filter(row => row.check === "api-latency");
	const detailRows = expandAssignedRows(baseRows);

	function toggleLocationMode() {
		const nextMode = locationMode === "manual" ? "auto" : "manual";

		setLocationMode(nextMode);

		if (nextMode === "auto") {
			setProbeLocation(probe.location);
		}
	}

	function toggleAsMode() {
		const nextMode = asMode === "manual" ? "auto" : "manual";

		setAsMode(nextMode);

		if (nextMode === "auto") {
			setProbeAsn(probe.asn);
		}
	}

	return (
		<section className={classNames(styles.card, floating && styles.floating)} aria-label="Probe detail">
			<div className={styles.header}>
				<span>Probe detail</span>
				<strong>
					{probe.name}
					<small> · uptime {probe.uptime}</small>
				</strong>
			</div>

			<div className={styles.fieldGrid}>
				<TextField className={styles.input} label="Probe name" value={probeName} onChange={event => setProbeName(event.currentTarget.value)} />
				<div className={styles.inputWithMode}>
					<TextField
						className={styles.input}
						label="Location (keywords search)"
						value={probeLocation}
						disabled={locationMode === "auto"}
						onChange={event => setProbeLocation(event.currentTarget.value)}
					/>
					<ModeToggle mode={locationMode} label="location detect mode" onClick={toggleLocationMode} />
				</div>
				<div className={styles.inputWithMode}>
					<TextField className={styles.input} label="AS" value={probeAsn} disabled={asMode === "auto"} onChange={event => setProbeAsn(event.currentTarget.value)} />
					<ModeToggle mode={asMode} label="AS detect mode" onClick={toggleAsMode} />
				</div>
			</div>

			<div className={styles.tableWrap}>
				<table className={styles.table}>
					<thead>
						<tr>
							<th>Assigned check</th>
							<th>Type</th>
							<th>Interval</th>
							<th>Jitter</th>
							<th>Latest</th>
						</tr>
					</thead>
					<tbody>
						{detailRows.map((row, index) => (
							<tr key={`${row.probe}-${row.check}-${index}`}>
								<td>{row.check}</td>
								<td>
									<span className={styles.checkType}>
										<span aria-hidden="true" />
										{row.type}
									</span>
								</td>
								<td>{row.interval}</td>
								<td>{row.jitter}</td>
								<td>{row.latest}</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</section>
	);
}

interface ModeToggleProps {
	mode: DetectionMode;
	label: string;
	onClick: () => void;
}

function ModeToggle({ mode, label, onClick }: ModeToggleProps) {
	const modeClass = mode === "manual" ? styles.modeButtonManual : styles.modeButtonAuto;

	return (
		<button className={classNames(styles.modeButton, modeClass)} type="button" aria-label={label} aria-pressed={mode === "auto"} onClick={onClick}>
			{mode}
		</button>
	);
}
