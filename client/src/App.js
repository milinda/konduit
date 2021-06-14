import {Col, Container, Navbar, Row} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css'
import './App.css';

function App() {
    return (
        <div>
            <Navbar bg="light">
                <Navbar.Brand href="#">Konduit</Navbar.Brand>
            </Navbar>
            <Container fluid="md">
                <Row>
                    <Col>Hello</Col>
                </Row>
            </Container>
        </div>
    );
}

export default App;
