import { Badge, Button, Terminal, TextField } from "@netstamp/ui";
import { type FormEvent, type MouseEvent, type KeyboardEvent as ReactKeyboardEvent, useEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { pathForRoute } from "../../../routes/routePaths";
import { NetworkMap } from "../../../shared/components/NetworkMap";
import { type CheckType, type Probe, type ProbeStatus, assignments, probes } from "../../../shared/utils/mockData";
import styles from "./ProductPages.module.css";

type ProbeView = "grid" | "map";
type ProbeSort = "heartbeat" | "name" | "asn";
type DetectionMode = "manual" | "auto";

interface AssignedRow {
	probe: string;
	check: string;
	type: CheckType;
	interval: string;
	jitter: string;
	latest: string;
}

const defaultProbeTags = ["Edge", "Home", "VPS", "Bare metal", "IPv6", "Web3", "Lab"];
const providerOptions = Array.from(new Set(probes.map(probe => probe.provider)));
const drawerCloseDurationMs = 180;
const installDetectionDurationMs = 1600;
const assignedOverflowRowCount = 18;
const createProbeSteps = [
	{ number: "01", title: "Name", copy: "Probe identity" },
	{ number: "02", title: "Install", copy: "Run command" },
	{ number: "03", title: "Details", copy: "Optional metadata" }
];

function asnNumber(asn: string) {
	return Number(asn.replace(/\D/g, "")) || 0;
}

function filterProbes(source: Probe[], search: string, statusFilter: "all" | ProbeStatus, providerFilter: string, sortKey: ProbeSort) {
	const term = search.trim().toLowerCase();
	const filtered = source.filter(probe => {
		const searchable = [probe.name, probe.location, probe.publicIp, probe.asn, probe.provider, probe.region, ...probe.tags].join(" ").toLowerCase();

		return (!term || searchable.includes(term)) && (statusFilter === "all" || probe.status === statusFilter) && (providerFilter === "all" || probe.provider === providerFilter);
	});

	if (sortKey === "name") {
		return filtered.sort((left, right) => left.name.localeCompare(right.name));
	}

	if (sortKey === "asn") {
		return filtered.sort((left, right) => asnNumber(left.asn) - asnNumber(right.asn));
	}

	return filtered;
}

function expandAssignedRows(rows: AssignedRow[]) {
	if (!rows.length) {
		return [];
	}

	return Array.from({ length: assignedOverflowRowCount }, (_, index) => {
		const row = rows[index % rows.length];
		const suffix = String(index + 1).padStart(2, "0");

		return {
			...row,
			check: index < rows.length ? row.check : `${row.check}-${suffix}`
		};
	});
}

export function ProbesPage() {
	const location = useLocation();
	const [view, setView] = useState<ProbeView>("grid");
	const [selectedId, setSelectedId] = useState("ams-edge-01");
	const [search, setSearch] = useState("");
	const [statusFilter, setStatusFilter] = useState<"all" | ProbeStatus>("all");
	const [providerFilter, setProviderFilter] = useState("all");
	const [sortKey, setSortKey] = useState<ProbeSort>("heartbeat");
	const wizardOpen = location.hash === "#new-probe";
	const selectedProbe = probes.find(probe => probe.id === selectedId) || probes[0];
	const visibleProbes = filterProbes(probes, search, statusFilter, providerFilter, sortKey);
	const assignedRows: AssignedRow[] = assignments.map(([probe, check, type, interval, jitter, latest]) => ({
		probe,
		check,
		type,
		interval,
		jitter,
		latest
	}));

	return (
		<section className={`${styles.probesScreen} ${view === "map" ? styles.probesScreenMap : ""}`}>
			{view === "grid" ? (
				<>
					<ProbePageHeader view={view} onViewChange={setView} />
					<div className={styles.probeGridLayout}>
						<ProbeList
							probes={visibleProbes}
							selectedId={selectedId}
							search={search}
							statusFilter={statusFilter}
							providerFilter={providerFilter}
							sortKey={sortKey}
							onSearchChange={setSearch}
							onStatusChange={setStatusFilter}
							onProviderChange={setProviderFilter}
							onSortChange={setSortKey}
							onSelect={setSelectedId}
						/>
						<div className={styles.probeLowerGrid}>
							<NetworkMap probes={probes} selectedId={selectedId} onSelect={setSelectedId} mode="detail" className={styles.probeMiniMap} />
							<ProbeDetail key={selectedProbe.id} probe={selectedProbe} assignedRows={assignedRows} />
						</div>
					</div>
				</>
			) : (
				<div className={styles.probeMapView}>
					<NetworkMap probes={probes} selectedId={selectedId} onSelect={setSelectedId} mode="fleet" className={styles.probeFullMap} />
					<ProbePageHeader view={view} onViewChange={setView} overlay />
					<ProbeDetail key={selectedProbe.id} probe={selectedProbe} assignedRows={assignedRows} floating />
				</div>
			)}

			{wizardOpen ? <NewProbeDrawer /> : null}
		</section>
	);
}

interface ProbePageHeaderProps {
	view: ProbeView;
	onViewChange: (view: ProbeView) => void;
	overlay?: boolean;
}

function ProbePageHeader({ view, onViewChange, overlay = false }: ProbePageHeaderProps) {
	return (
		<header className={[styles.probeHeader, overlay ? styles.probeHeaderOverlay : ""].filter(Boolean).join(" ")}>
			<div>
				<span className={[styles.probeKicker, overlay ? styles.probeKickerAccent : styles.probeKickerNeutral].join(" ")}>{overlay ? "Probe management" : "Last 24 hours"}</span>
				<h1>Probe Fleet</h1>
			</div>
			<div className={styles.probeHeaderActions}>
				<Button type="button" size="sm" variant={view === "grid" ? "secondary" : "ghost"} onClick={() => onViewChange("grid")}>
					Grid View
				</Button>
				<Button type="button" size="sm" variant={view === "map" ? "secondary" : "ghost"} onClick={() => onViewChange("map")}>
					Map View
				</Button>
				<Button className={styles.createProbeButton} size="sm" asChild>
					<Link to={`${pathForRoute("probes")}#new-probe`}>Create Probe</Link>
				</Button>
			</div>
		</header>
	);
}

interface ProbeListProps {
	probes: Probe[];
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

function ProbeList({ probes, selectedId, search, statusFilter, providerFilter, sortKey, onSearchChange, onStatusChange, onProviderChange, onSortChange, onSelect }: ProbeListProps) {
	function handleRowKeyDown(event: ReactKeyboardEvent<HTMLTableRowElement>, probeId: string) {
		if (event.key === "Enter" || event.key === " ") {
			event.preventDefault();
			onSelect(probeId);
		}
	}

	return (
		<section className={styles.probeListPanel} aria-label="Probe list">
			<div className={styles.probeListToolbar}>
				<span className={styles.probeListTitle}>Probe list</span>
				<input className={styles.toolbarControl} aria-label="Search probes" placeholder="Search" value={search} onChange={event => onSearchChange(event.currentTarget.value)} />
				<select className={styles.toolbarControl} aria-label="Filter status" value={statusFilter} onChange={event => onStatusChange(event.currentTarget.value as "all" | ProbeStatus)}>
					<option value="all">Status</option>
					<option value="Online">Online</option>
					<option value="Draining">Draining</option>
					<option value="Offline">Offline</option>
				</select>
				<select className={styles.toolbarControl} aria-label="Filter provider" value={providerFilter} onChange={event => onProviderChange(event.currentTarget.value)}>
					<option value="all">Provider</option>
					{providerOptions.map(provider => (
						<option key={provider} value={provider}>
							{provider}
						</option>
					))}
				</select>
				<select className={`${styles.toolbarControl} ${styles.sortControl}`} aria-label="Sort probes" value={sortKey} onChange={event => onSortChange(event.currentTarget.value as ProbeSort)}>
					<option value="heartbeat">Sort: Last Heartbeat</option>
					<option value="name">Sort: Probe Name</option>
					<option value="asn">Sort: AS</option>
				</select>
			</div>

			<div className={styles.probeTableWrap}>
				<table className={styles.probeTable}>
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
									className={probe.id === selectedId ? styles.selectedProbeRow : undefined}
									tabIndex={0}
									onClick={() => onSelect(probe.id)}
									onKeyDown={event => handleRowKeyDown(event, probe.id)}
								>
									<td>{probe.name}</td>
									<td>
										<span className={[styles.statusPill, styles[`status${probe.status}`]].filter(Boolean).join(" ")}>
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
										<span className={styles.probeTagList}>
											{probe.tags.map(tag => (
												<span className={styles.probeTag} key={tag}>
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
								<td className={styles.emptyProbeRow} colSpan={9}>
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

function NewProbeDrawer() {
	const navigate = useNavigate();
	const closeTimeoutRef = useRef<number | null>(null);
	const detectTimeoutRef = useRef<number | null>(null);
	const [closing, setClosing] = useState(false);
	const [currentStep, setCurrentStep] = useState(0);
	const [installStatus, setInstallStatus] = useState<"idle" | "detecting" | "detected">("idle");
	const [probeName, setProbeName] = useState("");
	const [probeLocation, setProbeLocation] = useState("");
	const [asn, setAsn] = useState("");
	const [tagOptions, setTagOptions] = useState(defaultProbeTags);
	const [selectedTags, setSelectedTags] = useState(["Edge"]);
	const [newTag, setNewTag] = useState("");
	const canCreate = probeName.trim().length > 0;
	const token = "NSTP_yoru-first-probe";
	const installCommand = [
		`sudo netstamp register --controller https://controller.netstamp.io --token ${token} --name "${probeName.trim() || "your-probe"}"`,
		"sudo systemctl enable --now netstamp-probe"
	].join("\n");

	useEffect(() => {
		return () => {
			if (closeTimeoutRef.current) {
				window.clearTimeout(closeTimeoutRef.current);
			}

			if (detectTimeoutRef.current) {
				window.clearTimeout(detectTimeoutRef.current);
			}
		};
	}, []);

	useEffect(() => {
		if (currentStep !== 1 || installStatus !== "detecting") {
			return undefined;
		}

		detectTimeoutRef.current = window.setTimeout(() => {
			setInstallStatus("detected");
			detectTimeoutRef.current = null;
		}, installDetectionDurationMs);

		return () => {
			if (detectTimeoutRef.current) {
				window.clearTimeout(detectTimeoutRef.current);
				detectTimeoutRef.current = null;
			}
		};
	}, [currentStep, installStatus]);

	useEffect(() => {
		function handleKeyDown(event: KeyboardEvent) {
			if (event.key !== "Escape" || closing || closeTimeoutRef.current) {
				return;
			}

			setClosing(true);
			closeTimeoutRef.current = window.setTimeout(() => navigate(pathForRoute("probes")), drawerCloseDurationMs);
		}

		window.addEventListener("keydown", handleKeyDown);
		return () => window.removeEventListener("keydown", handleKeyDown);
	}, [closing, navigate]);

	function closeDrawer() {
		if (closing || closeTimeoutRef.current) {
			return;
		}

		setClosing(true);
		closeTimeoutRef.current = window.setTimeout(() => navigate(pathForRoute("probes")), drawerCloseDurationMs);
	}

	function updateProbeName(value: string) {
		setProbeName(value);
		setInstallStatus("idle");
	}

	function startInstallDetection() {
		setInstallStatus("detecting");
		setCurrentStep(1);
	}

	function handleNameSubmit(event: FormEvent<HTMLFormElement>) {
		event.preventDefault();

		if (canCreate) {
			startInstallDetection();
		}
	}

	function toggleTag(tag: string) {
		setSelectedTags(current => (current.includes(tag) ? current.filter(value => value !== tag) : [...current, tag]));
	}

	function addTag() {
		const trimmedTag = newTag.trim();

		if (!trimmedTag) {
			return;
		}

		setTagOptions(current => (current.includes(trimmedTag) ? current : [...current, trimmedTag]));
		setSelectedTags(current => (current.includes(trimmedTag) ? current : [...current, trimmedTag]));
		setNewTag("");
	}

	function handleBackdropClick(event: MouseEvent<HTMLDivElement>) {
		if (event.target === event.currentTarget) {
			closeDrawer();
		}
	}

	return (
		<div className={`${styles.drawerBackdrop} ${closing ? styles.drawerBackdropClosing : ""}`} onClick={handleBackdropClick}>
			<aside className={`${styles.probeWizardDrawer} ${closing ? styles.probeWizardDrawerClosing : ""}`} aria-label="New probe wizard">
				<div className={styles.drawerHeader}>
					<div>
						<Badge tone="accent">New probe wizard</Badge>
						<h2>Create probe</h2>
						<p>Name the probe, install it on a host, then optionally annotate it with network metadata.</p>
					</div>
					<Button type="button" variant="ghost" size="sm" onClick={closeDrawer}>
						Close
					</Button>
				</div>

				<ol className={styles.stepTimeline} aria-label="Create probe progress">
					{createProbeSteps.map((step, index) => (
						<li className={`${styles.stepItem} ${index === currentStep ? styles.stepActive : ""} ${index < currentStep ? styles.stepComplete : ""}`} key={step.number}>
							<span>{step.number}</span>
							<strong>{step.title}</strong>
							<small>{step.copy}</small>
						</li>
					))}
				</ol>

				<div className={styles.workflowViewport}>
					<div className={styles.workflowTrack} style={{ transform: `translateX(-${currentStep * 100}%)` }}>
						<form className={styles.workflowPanel} aria-hidden={currentStep !== 0} onSubmit={handleNameSubmit}>
							<div className={styles.stepCopy}>
								<Badge tone="accent">Step 01</Badge>
								<h3>Enter probe name</h3>
								<p>This name is embedded in the registration command and shown in the probe fleet.</p>
							</div>

							<TextField label="Probe name" value={probeName} placeholder="taipei-home-01" required disabled={currentStep !== 0} onChange={event => updateProbeName(event.currentTarget.value)} />

							<div className={styles.drawerActions}>
								<Button type="submit" disabled={!canCreate || currentStep !== 0}>
									Continue to install
								</Button>
								<p className={styles.drawerHint}>Use a stable hostname-style label so results are easy to scan later.</p>
							</div>
						</form>

						<section className={styles.workflowPanel} aria-hidden={currentStep !== 1}>
							<div className={styles.stepCopy}>
								<Badge tone={installStatus === "detected" ? "success" : "warning"}>{installStatus === "detected" ? "Probe detected" : "Auto detecting"}</Badge>
								<h3>Install the probe</h3>
								<p>Run this command on the host. The wizard watches for the first heartbeat and unlocks metadata when the probe is detected.</p>
							</div>

							<div className={styles.registrationBlock}>
								<div className={styles.tokenLine}>
									<span>Registration token</span>
									<strong>{token}</strong>
								</div>
								<Terminal title="install command" meta="copy to host">
									{installCommand}
								</Terminal>
							</div>

							<div className={styles.detectCard}>
								<Badge tone={installStatus === "detected" ? "success" : "warning"}>{installStatus === "detected" ? "Heartbeat received" : "Listening for heartbeat"}</Badge>
								<strong>{installStatus === "detected" ? `${probeName.trim()} is online` : "Waiting for install to finish"}</strong>
								<p>{installStatus === "detected" ? "The controller accepted the first signed result stream." : "This frontend mock auto-detects shortly after the install step starts."}</p>
							</div>

							<div className={styles.drawerActions}>
								<Button type="button" variant="ghost" disabled={currentStep !== 1} onClick={() => setCurrentStep(0)}>
									Back
								</Button>
								<Button type="button" disabled={installStatus !== "detected" || currentStep !== 1} onClick={() => setCurrentStep(2)}>
									Continue to details
								</Button>
							</div>
						</section>

						<section className={styles.workflowPanel} aria-hidden={currentStep !== 2}>
							<div className={styles.stepCopy}>
								<Badge tone="accent">Step 03</Badge>
								<h3>Add optional metadata</h3>
								<p>Customize location, AS, and tags or we can guess it by ourselves.</p>
							</div>

							<div className={styles.formGrid}>
								<TextField
									label="Location (optional)"
									value={probeLocation}
									placeholder="Taipei, Taiwan"
									disabled={currentStep !== 2}
									onChange={event => setProbeLocation(event.currentTarget.value)}
								/>
								<TextField label="AS (optional)" value={asn} placeholder="AS3462" disabled={currentStep !== 2} onChange={event => setAsn(event.currentTarget.value)} />
							</div>

							<div className={styles.tagPicker}>
								<span className={styles.drawerFieldLabel}>Tags</span>
								<div className={styles.tagCloud}>
									{tagOptions.map(tag => (
										<button
											className={`${styles.tagButton} ${selectedTags.includes(tag) ? styles.tagSelected : ""}`}
											key={tag}
											type="button"
											disabled={currentStep !== 2}
											onClick={() => toggleTag(tag)}
										>
											{tag}
										</button>
									))}
								</div>
								<div className={styles.tagCreate}>
									<TextField
										label="Create tag"
										value={newTag}
										placeholder="backbone"
										disabled={currentStep !== 2}
										onChange={event => setNewTag(event.currentTarget.value)}
										onKeyDown={event => {
											if (event.key === "Enter") {
												event.preventDefault();
												addTag();
											}
										}}
									/>
									<Button type="button" variant="outline" disabled={currentStep !== 2} onClick={addTag}>
										Add tag
									</Button>
								</div>
							</div>

							<div className={styles.drawerActions}>
								<Button type="button" variant="ghost" disabled={currentStep !== 2} onClick={() => setCurrentStep(1)}>
									Back
								</Button>
								<Button type="button" disabled={currentStep !== 2} onClick={closeDrawer}>
									Save details
								</Button>
								<Button type="button" variant="outline" disabled={currentStep !== 2} onClick={closeDrawer}>
									Skip
								</Button>
							</div>
						</section>
					</div>
				</div>
			</aside>
		</div>
	);
}

interface ProbeDetailProps {
	probe: Probe;
	assignedRows: AssignedRow[];
	floating?: boolean;
}

function ProbeDetail({ probe, assignedRows, floating = false }: ProbeDetailProps) {
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
		<section className={[styles.probeDetailCard, floating ? styles.probeDetailFloating : ""].filter(Boolean).join(" ")} aria-label="Probe detail">
			<div className={styles.probeDetailHeader}>
				<span>Probe detail</span>
				<strong>
					{probe.name}
					<small> · uptime {probe.uptime}</small>
				</strong>
			</div>

			<div className={styles.probeFieldGrid}>
				<TextField className={styles.probeDetailInput} label="Probe name" value={probeName} onChange={event => setProbeName(event.currentTarget.value)} />
				<div className={styles.probeInputWithMode}>
					<TextField
						className={styles.probeDetailInput}
						label="Location (keywords search)"
						value={probeLocation}
						disabled={locationMode === "auto"}
						onChange={event => setProbeLocation(event.currentTarget.value)}
					/>
					<ModeToggle mode={locationMode} label="location detect mode" onClick={toggleLocationMode} />
				</div>
				<div className={styles.probeInputWithMode}>
					<TextField className={styles.probeDetailInput} label="AS" value={probeAsn} disabled={asMode === "auto"} onChange={event => setProbeAsn(event.currentTarget.value)} />
					<ModeToggle mode={asMode} label="AS detect mode" onClick={toggleAsMode} />
				</div>
			</div>

			<div className={styles.assignedTableWrap}>
				<table className={styles.assignedTable}>
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
	const modeClass = mode === "manual" ? styles.fieldModeButtonManual : styles.fieldModeButtonAuto;

	return (
		<button className={[styles.fieldModeButton, modeClass].join(" ")} type="button" aria-label={label} aria-pressed={mode === "auto"} onClick={onClick}>
			{mode}
		</button>
	);
}
