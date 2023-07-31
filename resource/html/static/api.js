const API_BASE = '/api'
const RECORD_API_BASE = API_BASE + '/record'
const SYSTEM_API_BASE = API_BASE + '/system'

Vue.use({
    install: function (Vue, _) {
        const http = axios.create({
            timeout: 5000,
        })

        http.interceptors.response.use(function (response) {
            const data = response.data
            if (!data) {
                Vue.prototype.$message.error('接口请求失败，请重试')
                return undefined
            }
            if (data.errmsg) {
                Vue.prototype.$message.error(data.errmsg)
            }
            if (data.data !== undefined) {
                return data.data
            }
            return data.errno === 0
        }, function (error) {
            Vue.prototype.$message.error(error.message)
            return undefined
        })

        Vue.prototype.$api = {
            record: {
                async list() {
                    return await http.get(RECORD_API_BASE + '/')
                },
                async get(domain) {
                    return await http.get(RECORD_API_BASE + '/' + encodeURIComponent(domain))
                },
                async set(domain, ipv4, ipv6) {
                    return await http.put(RECORD_API_BASE + '/', {
                        domain: domain,
                        ipv4: ipv4,
                        ipv6: ipv6,
                        force_update_ipv4: true,
                        force_update_ipv6: true,
                    })
                },
                async delete(domain) {
                    return await http.delete(RECORD_API_BASE + '/' + encodeURIComponent(domain))
                }
            },
            system: {
                async get() {
                    return await http.get(SYSTEM_API_BASE + '/')
                },
                async set(web, dns) {
                    return await http.put(SYSTEM_API_BASE + '/', {
                        web: web,
                        dns: dns,
                    })
                },
            }
        }
    }
})
