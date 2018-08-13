import * as React from "react";
import * as ReactDOM from "react-dom";
import {Hello} from "./components/hello";

ReactDOM.render(
	<div className="container">
		<h1>MyEvents</h1>
		<Hello name="Sean"/>
	</div>,
	document.getElementById("myevents-app")
 );
