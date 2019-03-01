
// CTF{R34d1ng_src_c0d3_1s_my_f4v0rit3_h0bby}

login = new Vue({
    el: '.login',
    data: {
        username: '',
        password: ''
    },
    methods: {
        submit: function(e) {
            console.log('here we go!');
            console.log(this.username);
            console.log(this.password);
            $cookies.set('username', this.username);
            $cookies.set('password', this.password);
            location.reload();
        }
    }
});


