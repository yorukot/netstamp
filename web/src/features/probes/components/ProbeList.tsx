import type { KeyboardEvent as ReactKeyboardEvent } from "react";
import { classNames } from "../../../shared/utils/classNames";
import type { Probe, ProbeStatus } from "../../../shared/utils/mockData";
import styles from "./ProbeList.module.css";
import type { ProbeSort } from "./types";

interface ProbeListProps {
	probes: Probe[];
	providerOptions: string[];
	selectedId: string;
	search: string;
	statusFilter: "all" | ProbeStatus;
	providerFilter: string;
	sortKey: ProbeSort;
	onSearchChange: (value: string) => void;
	onStatusChange: (value: "all" | ProbeStatus) => void;
	onProviderChange: (value: string) => void;
	onSortChange: (value: ProbeSort) => void;
	onSelect: (probeId: string) => void;
}

export function ProbeList({
	probes,
	providerOptions,
	selectedId,
	search,
	statusFilter,
	providerFilter,
	sortKey,
	onSearchChange,
	onStatusChange,
	onProviderChange,
	onSortChange,
	onSelect
}: ProbeListProps) {
	function handleRowKeyDown(event: ReactKeyboardEvent<HTMLTableRowElement>, probeId: string) {
		if (event.key === "Enter" || event.key === " ") {
			event.preventDefault();
			onSelect(probeId);
		}
	}

	return (
		<section className={styles.panel} aria-label="Probe list">
			<div className={styles.toolbar}>
				<span className={styles.title}>Probe list</span>
				<input className={styles.control} aria-label="Search probes" placeholder="Search" value={search} onChange={event => onSearchChange(event.currentTarget.value)} />
				<select className={styles.control} aria-label="Filter status" value={statusFilter} onChange={event => onStatusChange(event.currentTarget.value as "all" | ProbeStatus)}>
					<option value="all">Status</option>
					<option value="Online">Online</option>
					<option value="Draining">Draining</option>
					<option value="Offline">Offline</option>
				</select>
				<select className={styles.control} aria-label="Filter provider" value={providerFilter} onChange={event => onProviderChange(event.currentTarget.value)}>
					<option value="all">Provider</option>
					{providerOptions.map(provider => (
						<option key={provider} value={provider}>
							{provider}
						</option>
					))}
				</select>
				<select className={classNames(styles.control, styles.sortControl)} aria-label="Sort probes" value={sortKey} onChange={event => onSortChange(event.currentTarget.value as ProbeSort)}>
					<option value="heartbeat">Sort: Last Heartbeat</option>
					<option value="name">Sort: Probe Name</option>
					<option value="asn">Sort: AS</option>
				</select>
			</div>

			<div className={styles.tableWrap}>
				<table className={styles.table}>
					<thead>
						<tr>
							<th>Probe name</th>
							<th>Status</th>
							<th>location</th>
							<th>Public IP</th>
							<th>AS</th>
							<th>Support IP Family</th>
							<th>last heartbeat</th>
							<th>tags</th>
							<th>Version</th>
						</tr>
					</thead>
					<tbody>
						{probes.length ? (
							probes.map(probe => (
								<tr
									key={probe.id}
									className={probe.id === selectedId ? styles.selectedRow : undefined}
									tabIndex={0}
									onClick={() => onSelect(probe.id)}
									onKeyDown={event => handleRowKeyDown(event, probe.id)}
								>
									<td>{probe.name}</td>
									<td>
										<span className={classNames(styles.statusPill, styles[`status${probe.status}` as keyof typeof styles])}>
											<span aria-hidden="true" />
											{probe.status}
										</span>
									</td>
									<td>{probe.location}</td>
									<td>{probe.publicIp}</td>
									<td>{probe.asn}</td>
									<td>{probe.ipFamily}</td>
									<td>{probe.lastHeartbeat}</td>
									<td>
										<span className={styles.tagList}>
											{probe.tags.map(tag => (
												<span className={styles.tag} key={tag}>
													{tag}
												</span>
											))}
										</span>
									</td>
									<td>{probe.version}</td>
								</tr>
							))
						) : (
							<tr>
								<td className={styles.emptyRow} colSpan={9}>
									No probes found
								</td>
							</tr>
						)}
					</tbody>
				</table>
			</div>
		</section>
	);
}
