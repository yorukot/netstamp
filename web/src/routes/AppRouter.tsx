import { createBrowserRouter, Navigate as RouterNavigate, RouterProvider, useNavigate } from "react-router-dom";
import { AuthPage } from "../features/auth/components/AuthPage";
import { OnboardingPage } from "../features/auth/components/OnboardingPage";
import { AlertsPage } from "../features/dashboard/components/AlertsPage";
import { ChecksPage } from "../features/dashboard/components/ChecksPage";
import { ComponentDemoPage } from "../features/dashboard/components/ComponentDemoPage";
import { DashboardPage } from "../features/dashboard/components/DashboardPage";
import { InsightPage } from "../features/dashboard/components/InsightPage";
import { LandingPage } from "../features/dashboard/components/LandingPage";
import { ProbesPage } from "../features/dashboard/components/ProbesPage";
import { SettingsPage } from "../features/dashboard/components/SettingsPage";
import { TeamPage } from "../features/dashboard/components/TeamPage";
import { AppShell } from "../layouts/AppShell";
import type { AppRoute, Navigate } from "../shared/utils/mockData";
import { pathForRoute } from "./routePaths";

function appRoutePath(route: AppRoute) {
	return pathForRoute(route).slice(1);
}

function useRouteNavigate(): Navigate {
	const navigate = useNavigate();

	return (route, hash) => navigate(`${pathForRoute(route)}${hash ?? ""}`);
}

function LandingRoute() {
	const navigate = useRouteNavigate();

	return <LandingPage navigate={navigate} />;
}

interface AuthRouteProps {
	mode: "login" | "register";
}

function AuthRoute({ mode }: AuthRouteProps) {
	const navigate = useRouteNavigate();

	return <AuthPage mode={mode} navigate={navigate} />;
}

function OnboardingRoute() {
	const navigate = useRouteNavigate();

	return <OnboardingPage navigate={navigate} />;
}

function DashboardRoute() {
	const navigate = useRouteNavigate();

	return <DashboardPage navigate={navigate} />;
}

const router = createBrowserRouter([
	{ path: pathForRoute("landing"), element: <LandingRoute /> },
	{ path: pathForRoute("login"), element: <AuthRoute mode="login" /> },
	{ path: pathForRoute("register"), element: <AuthRoute mode="register" /> },
	{ path: pathForRoute("onboarding"), element: <OnboardingRoute /> },
	{
		element: <AppShell />,
		children: [
			{ path: appRoutePath("dashboard"), element: <DashboardRoute /> },
			{ path: appRoutePath("probes"), element: <ProbesPage /> },
			{ path: appRoutePath("insight"), element: <InsightPage /> },
			{ path: appRoutePath("checks"), element: <ChecksPage /> },
			{ path: appRoutePath("alerts"), element: <AlertsPage /> },
			{ path: appRoutePath("team"), element: <TeamPage /> },
			{ path: appRoutePath("settings"), element: <SettingsPage /> },
			{ path: appRoutePath("components"), element: <ComponentDemoPage /> }
		]
	},
	{ path: "*", element: <RouterNavigate to={pathForRoute("landing")} replace /> }
]);

export function AppRouter() {
	return <RouterProvider router={router} />;
}
