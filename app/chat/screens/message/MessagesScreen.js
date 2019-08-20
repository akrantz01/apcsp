import React from 'react';
import {Component} from 'react';
import {RefreshControl, SectionList} from 'react-native';
import {Button} from 'react-native-elements';

export default class MessagesScreen extends Component {
    constructor(props) {
        super(props);
        this.state = {
            refreshing: false,
        };
    }

    static navigationOptions = {
        title: 'Messages',
    };

    _onRefresh() {
        this.setState({refreshing: true});
        // .then( () => {
        this.setState({refreshing: false});
        // });
    }

    render() {
        const {navigate} = this.props.navigation;
        return (
            <SectionList
                refreshControl={
                    <RefreshControl
                        refreshing={this.state.refreshing}
                        onRefresh={this._onRefresh.bind(this)}
                    />
                }
                sections={[
                    // homogeneous rendering between sections
                    {data: [{key: 'Jane'}], title: 'Jane'},
                    {data: [{key: 'Bob'}], title: 'Bob'},
                ]}
                renderItem={({item}) => (
                    <Button
                        title={item.key}
                        onPress={() => navigate('Thread', {name: item.key})}
                    />
                )}
                getItem={null}
                getItemCount={null}
            />
        );
    }
}
