import * as React from "react";
import * as ReactDOM from "react-dom";
import {EventListContainer} from "./components/event_list_container";

class App extends React.Component<{}, {}> {
	render() {
		return <div className="container">
			<h1>MyEvents</h1>
			<EventListContainer eventListURL="http://localhost:8181"/>
		</div>
	}
}

ReactDOM.render(
	<App/>,
	document.getElementById("myevents-app")
 );
