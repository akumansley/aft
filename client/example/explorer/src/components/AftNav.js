import React, { Component } from "react";
import {
  Dropdown,
  DropdownButton,
  Button,
  Nav,
  NavItem,
  Navbar
} from "react-bootstrap";
import { cap, isNonEmptyList } from "../util";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus } from "@fortawesome/free-solid-svg-icons";

class AftNav extends Component {
  constructor() {
    super();
    this.state = { name: "More", overflow: 3 };
    this.handleSelect = this.handleSelect.bind(this);
    this.handleNew = this.handleNew.bind(this);
  }

  overflow(list) {
    return (
      <DropdownButton
        variant="outline-primary"
        title={this.state.name}
        as={NavItem}
      >
        {list.map(value => (
          <Dropdown.Item
            eventKey={value}
            key={value}
            onSelect={e => this.handleSelect(e)}
          >
            {cap(value)}
          </Dropdown.Item>
        ))}
      </DropdownButton>
    );
  }

  nav(list, len) {
    var ovr = "";
    if (len < this.props.list.length) {
      ovr = this.overflow(this.props.list.slice(len));
    }
    return (
      <Nav defaultActiveKey={this.props.name}>
        {this.props.list.slice(0, len).map(value => (
          <Nav.Item key={value}>
            <Nav.Link eventKey={value} onSelect={e => this.handleSelect(e)}>
              {cap(value)}
            </Nav.Link>
          </Nav.Item>
        ))}
        {ovr}
      </Nav>
    );
  }

  handleSelect(e) {
    this.props.handleSelect(e);
    if (this.props.list.indexOf(e) >= this.state.overflow) {
      this.setState({ name: cap(e) });
    } else if (this.props.list.indexOf(e) !== -1) {
      this.setState({ name: "More" });
    }
  }

  handleNew(e) {
    this.props.handleNew(e);
  }

  render() {
    if (isNonEmptyList(this.props.list)) {
      return (
        <Navbar bg="light" expand="lg">
          <Navbar.Brand>Aft</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            {this.nav(this.props.list, this.state.overflow)}
            <Button
              variant="outline-secondary"
              as={NavItem}
              key={"new"}
              style={{ marginLeft: "10px" }}
              onClick={e => this.handleNew(e)}
            >
              <FontAwesomeIcon icon={faPlus} /> New
            </Button>
          </Navbar.Collapse>
        </Navbar>
      );
    }
    return <div></div>;
  }
}

export default AftNav;
