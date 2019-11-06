import React from 'react';
import {Component} from 'react';
import {
    RefreshControl,
    SectionList,
    StyleSheet,
    Text,
    View,
    TouchableOpacity,
    StatusBar,
    SafeAreaView,
} from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import {Avatar, Icon} from 'react-native-elements';
import Dialog from 'react-native-dialog';

export default class Messages extends Component {
    constructor(props) {
        super(props);
        this.state = {
            refreshing: false,
            edit: false,
            dialogVisible: false,
            net: false,
        };
    }

    edit = false;

    static navigationOptions = {
        title: 'Messages',
        header: null,
    };

    _onRefresh() {
        this.setState({refreshing: true});

        setTimeout(() => {
            this.getChatsOverNetwork();
            this.setState({refreshing: false});
        }, 1000);
    }

    getData() {
        return [
            {data: [{key: 0, name: 'Alex Beaver', unread: true, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 1, name: 'Alex Krantz', unread: true, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 2, name: 'Aidan Sacco', unread: false, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 3, name: 'Guy Wilks', unread: true, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 4, name: 'Kai Frondsal', unread: false, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 5, name: 'Eric Ettlin', unread: false, avatar: 'https://placeimg.com/140/140/any'}]},
            {data: [{key: 6, name: 'Dylan Pratt', unread: false, avatar: 'https://placeimg.com/140/140/any'}]},
        ];
    }

    getChatsOverNetwork() {
        this.getData();
    }

    deleteThread() {
        this.setState({dialogVisible: false});
        this.toggleEdit();
    }

    toggleEdit() {
        let b = !this.state.edit;
        this.setState({edit: b});
    }

    render() {
        const {navigate} = this.props.navigation;
        return (
            <LinearGradient style={styles.background} colors={['#00af3a', '#005baf']}>
                <Dialog.Container visible={this.state.dialogVisible}>
                    <Dialog.Title>Delete Thread?</Dialog.Title>
                    <Dialog.Description>
                        Do you want to delete this thread? You cannot undo this action.
                    </Dialog.Description>
                    <Dialog.Button label="Cancel" bold onPress={() => this.setState({dialogVisible: false})} />
                    <Dialog.Button label="Delete" color={'red'} onPress={() => this.deleteThread()} />
                </Dialog.Container>

                <StatusBar barStyle={'light-content'} networkActivityIndicatorVisible={this.state.net} />
                <SafeAreaView style={styles.safeArea} maxWidth={600}>
                    <View style={styles.bar}>
                        <Text style={styles.text}>Messages</Text>
                        <TouchableOpacity onPress={() => navigate('Settings')}>
                            <View style={styles.settingsButton}>
                                <Icon name={'settings'} type={'feather'} color={'black'} />
                            </View>
                        </TouchableOpacity>
                    </View>
                    <View style={styles.buttonBar}>
                        <View style={styles.buttonContainerMessage}>
                            <TouchableOpacity
                                onPress={() =>
                                    navigate({
                                        routeName: 'New',
                                        params: {
                                            transition: 'transition',
                                        },
                                    })
                                }>
                                <View style={styles.buttonBarButtonBackground}>
                                    <View style={styles.buttonBarContents}>
                                        <Icon name={'plus'} type={'feather'} color={'white'} />
                                        <Text style={styles.buttonText}>New Message</Text>
                                    </View>
                                </View>
                            </TouchableOpacity>
                        </View>
                        <View style={styles.buttonContainerEdit}>
                            <TouchableOpacity onPress={() => this.toggleEdit()}>
                                <View
                                    style={[
                                        styles.buttonBarButtonBackground,
                                        this.state.edit ? styles.editSelected : styles.editButton,
                                    ]}>
                                    <Text style={this.state.edit ? styles.editTextSelected : styles.editText}>
                                        Edit
                                    </Text>
                                </View>
                            </TouchableOpacity>
                        </View>
                    </View>
                    <SectionList
                        refreshControl={
                            <RefreshControl
                                refreshing={this.state.refreshing}
                                onRefresh={this._onRefresh.bind(this)}
                                tintColor={'white'}
                            />
                        }
                        sections={this.getData()}
                        contentContainerStyle={styles.sectionContainer}
                        renderItem={({item}) => (
                            <TouchableOpacity
                                onPress={() =>
                                    this.state.edit
                                        ? null
                                        : navigate('Thread', {
                                              name: item.name, //.split(' ')[0],
                                              chatID: item.key,
                                          })
                                }>
                                <View style={styles.card}>
                                    <LinearGradient
                                        start={{x: 0, y: 0}}
                                        end={{x: 1, y: 0}}
                                        colors={['white', '#dddddd']}
                                        style={styles.cardGrad}>
                                        <Avatar
                                            rounded
                                            title={item.name[0]}
                                            source={{
                                                uri: item.avatar,
                                            }}
                                            size={50}
                                        />
                                        <View style={styles.cardTextContainer}>
                                            <Text style={styles.name} numberOfLines={1}>
                                                {item.name}
                                            </Text>
                                            <Text style={styles.cardText} numberOfLines={1}>
                                                this is the previous message
                                            </Text>
                                        </View>
                                        {this.state.edit ? (
                                            <TouchableOpacity
                                                style={styles.deleteButton}
                                                onPress={() =>
                                                    this.setState({
                                                        dialogVisible: true,
                                                    })
                                                }>
                                                <Icon name={'x'} type={'feather'} color={'red'} size={20} />
                                            </TouchableOpacity>
                                        ) : null}
                                    </LinearGradient>
                                </View>
                            </TouchableOpacity>
                        )}
                        getItem={null}
                        getItemCount={null}
                    />
                </SafeAreaView>
            </LinearGradient>
        );
    }
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        paddingTop: '15%',
        paddingBottom: 0,
        height: '100%',
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'center',
    },
    bar: {
        display: 'flex',
        width: '93%',
        flexDirection: 'row',
        alignItems: 'center',
    },
    settingsButton: {
        flex: 1,
        maxHeight: 42,
        backgroundColor: 'white',
        borderRadius: 30,
        paddingTop: 8,
        paddingBottom: 3,
        paddingLeft: 19,
        paddingRight: 19,
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
    editButton: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        padding: 0,
        backgroundColor: '#5219d8',
    },
    editSelected: {
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        padding: 0,
        backgroundColor: 'red',
    },
    editText: {
        color: 'white',
        fontWeight: 'normal',
    },
    editTextSelected: {
        color: 'white',
        fontWeight: 'bold',
    },
    buttonText: {
        marginLeft: 5,
        color: 'white',
    },
    buttonBar: {
        height: 45,
        padding: 0,
        marginLeft: 20,
        marginRight: 20,
        marginBottom: 10,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        display: 'flex',
    },
    buttonBarButtonBackground: {
        height: '100%',
        borderRadius: 40,
        backgroundColor: '#5219d8',
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
    buttonContainerMessage: {
        height: '100%',
        borderRadius: 40,
        marginLeft: 0,
        marginRight: 5,
        flex: 4,
    },
    buttonContainerEdit: {
        height: '100%',
        borderRadius: 40,
        marginLeft: 0,
        marginRight: 5,
        flex: 1,
    },
    buttonBarContents: {
        padding: 10,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
    },
    cardGrad: {
        padding: 20,
        borderRadius: 30,
        justifyContent: 'center',
    },
    avatar: {
        position: 'absolute',
        left: 0,
    },
    text: {
        flex: 1,
        marginLeft: 20,
        marginBottom: 10,
        fontSize: 40,
        color: 'white',
        fontWeight: '800',
    },
    cardTextContainer: {
        padding: 20,
        position: 'absolute',
        right: 0,
        top: 0,
        bottom: 0,
        left: 70,
    },
    cardText: {
        marginTop: 5,
        color: '#666666',
        fontSize: 16,
    },
    name: {
        fontSize: 20,
        fontWeight: '600',
    },
    card: {
        marginLeft: 20,
        marginRight: 20,
        marginBottom: 10,
        marginTop: 10,
        borderRadius: 30,
        backgroundColor: 'white',
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
    safeArea: {
        height: '108%',
        flex: 1,
    },
    sectionContainer: {
        paddingBottom: 10,
    },
    deleteButton: {
        position: 'absolute',
        right: 15,
        top: 15,
    },
});
