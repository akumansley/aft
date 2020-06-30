import React, { Component } from "react";
import client from "../client";
import Form from "@rjsf/core";
import { isObject } from "../util";

class AftForm extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.handleSubmit = this.handleSubmit.bind(this);
    if (this.props.name !== "") {
      let load = client.views.rpc.reactForm({
        data: { model: this.props.name }
      });
      load.then(obj => {
        if (obj !== undefined) {
          this.setState({ schema: obj });
        }
      });
    }
  }

  handleSubmit(e) {
    this.setState(prevState => {
      return {
        ...prevState,
        formData: e.formData
      };
    });
    let load = client.api[this.props.name].create({ data: e.formData });
    let test = client.views.rpc.validate({data: { schema: this.state.schema, data: e.formData }});
    load.then(
      obj => {
        this.props.handleSubmit(this.props.name);
      },
      err => {
        this.props.handleError(err);
      }
    );
  }

  render() {
    if (isObject(this.state.schema)) {
      return (
        <Form
          schema={this.state.schema}
          formData={this.state.formData}
          onSubmit={fd => this.handleSubmit(fd)}
        />
      );
    }
    return <div></div>;
  }
}

export default AftForm;
