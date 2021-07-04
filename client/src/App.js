import React from 'react';
import {Col, Container, Navbar, Row} from "react-bootstrap";
import {w3cwebsocket as W3CWebSocket} from 'websocket';
import 'bootstrap/dist/css/bootstrap.min.css'
import './App.css';

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            counter: "0"
        }
        this.onWebSocketOpen = this.onWebSocketOpen.bind(this);
        this.onWebSocketMessage = this.onWebSocketMessage.bind(this);
    }

    onWebSocketOpen(event) {
        console.log("Websocket connected: " + event);
    }

    onWebSocketMessage(event) {
        this.setState({
            counter: event.data
        })
    }

    componentDidMount() {
        this.wsClient = new W3CWebSocket(`ws://${window.location.host}/ws`);
        this.wsClient.onopen = this.onWebSocketOpen;
        this.wsClient.onmessage = this.onWebSocketMessage;
    }

    render() {
        return (
            <div>
                <Navbar bg="light">
                    <Navbar.Brand href="#">Konduit</Navbar.Brand>
                </Navbar>
                <Container fluid="md">
                    <Row>
                        <Col><p>{this.state.counter}</p></Col>
                    </Row>
                </Container>
            </div>
        );
    }
}

export default App;
