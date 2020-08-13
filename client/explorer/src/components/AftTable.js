import React, { Component } from "react";
import aft from "../aft";
import { Table } from "react-bootstrap";
import { cap, isNonEmptyList, isObject } from "../util";

class AftTable extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.refresh = this.refresh.bind(this);
    this.refresh();
  }

  refresh() {
    if (this.props.name !== "") {
      let load = aft.api[this.props.name].findMany({ where: {} });
      load.then(obj => {
        if (obj !== undefined) {
          this.setState({ list: obj });
        }
      });
    }
  }

  componentDidUpdate(prevProps) {
    if (prevProps.name !== this.props.name) {
      this.refresh();
    }
  }

  title(elemn) {
    if (isObject(elemn)) {
      var count = 0;
      return (
        <tr key={count++}>
          {Object.keys(elemn).map(value => (
            <th scope="col" key={count++}>
              {cap(value)}
            </th>
          ))}
        </tr>
      );
    }
    return <tr></tr>;
  }

  row(elemn) {
    if (isObject(elemn)) {
      var count = 0;
      return (
        <tr key={elemn.id}>
          {Object.values(elemn).map(value => (
            <td key={elemn.id.concat(count++)}>{value.toString()}</td>
          ))}
        </tr>
      );
    }
    return <tr></tr>;
  }

  render() {
    if (isNonEmptyList(this.state.list)) {
      const title = <thead>{this.title(this.state.list[0])}</thead>;
      const rows = (
        <tbody>
          {Object.values(this.state.list).map(elemn => this.row(elemn))}
        </tbody>
      );
      return (
        <Table
          className="table table-striped table-bordered table-hover"
          size="sm"
          responsive
        >
          {title}
          {rows}
        </Table>
      );
    }

    return <div>{cap(this.props.name)} model has no records</div>;
  }
}

export default AftTable;
