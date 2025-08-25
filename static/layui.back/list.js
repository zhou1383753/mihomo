layui.use(['form','layer','laydate','table','laytpl'],function(){
    var form = layui.form,
        layer = parent.layer === undefined ? layui.layer : top.layer,
        $ = layui.jquery,
        laydate = layui.laydate,
        laytpl = layui.laytpl,
        table = layui.table;

    //新闻列表
    var tableIns = table.render({
        elem: '#newsList',
        url : '/index/user/zhlist',
        cellMinWidth : 95,
        page : true,
        limit : 15,
        limits : [10,15,20,25],
        id : "newsListTable",
        cols : [[
            {type: 'checkbox', fixed: 'left'}

            ,{field: 'user', title: '账号', width:120}
            ,{field: 'usetime', title: '密码', width: 80}
            ,{field: 'value', title: '连接数', width: 60}
            ,{field: 'kmtype_id', title: '类目', width: 120, sort: true, totalRow: true,templet:"#newsStatus"}
            ,{field: 'kmtype_lb', title: '类别', width:80,templet:"#newsStatuss"}
            ,{field: 'addtime', title: '开卡时间', width: 180, sort: true, totalRow: true}
            ,{field: 'edntime', title: '到期时间', width: 180, sort: true, totalRow: true}
            ,{field: 'money', title: '扣费', width: 80, sort: true, totalRow: true}
            ,{field: 'km', title: '备注', width: 180, sort: true, totalRow: true}
            ,{fixed: 'right', width: 165, align:'center', toolbar: '#barDemo'}
        ]]
    });

    //工具栏事件
    table.on('toolbar(test)', function(obj){
        var checkStatus = table.checkStatus(obj.config.id);
        switch(obj.event){
            case 'getCheckData':
                var data = checkStatus.data;
                layer.alert(JSON.stringify(data));
                break;
            case 'getCheckLength':
                var data = checkStatus.data;
                layer.msg('选中了：'+ data.length + ' 个');
                break;
            case 'isAll':
                layer.msg(checkStatus.isAll ? '全选': '未全选')
                break;
        };
    });



    //是否置顶
    form.on('switch(newsTop)', function(data){
        var index = layer.msg('修改中，请稍候',{icon: 16,time:false,shade:0.8});
        setTimeout(function(){
            layer.close(index);
            if(data.elem.checked){
                layer.msg("置顶成功！");
            }else{
                layer.msg("取消置顶成功！");
            }
        },500);
    })

    //搜索【此功能需要后台配合，所以暂时没有动态效果演示】
    $(".search_btn").on("click",function(){
        if($(".searchVal").val() != ''){
            table.reload("newsListTable",{
                page: {
                    curr: 1 //重新从第 1 页开始
                },
                where: {
                    key: $(".searchVal").val()  //搜索的关键字
                }
            })
        }else{
            layer.msg("请输入搜索的内容");
        }
    });

    //添加文章
    function addNews(edit){
        var index = layui.layer.open({
            title : "添加文章",
            type : 2,
            content : "newsAdd.html",
            success : function(layero, index){
                var body = layui.layer.getChildFrame('body', index);
                if(edit){
                    body.find(".newsName").val(edit.newsName);
                    body.find(".abstract").val(edit.abstract);
                    body.find(".thumbImg").attr("src",edit.newsImg);
                    body.find("#news_content").val(edit.content);
                    body.find(".newsStatus select").val(edit.newsStatus);
                    body.find(".openness input[name='openness'][title='"+edit.newsLook+"']").prop("checked","checked");
                    body.find(".newsTop input[name='newsTop']").prop("checked",edit.newsTop);
                    form.render();
                }
                setTimeout(function(){
                    layui.layer.tips('点击此处返回文章列表', '.layui-layer-setwin .layui-layer-close', {
                        tips: 3
                    });
                },500)
            }
        })
        layui.layer.full(index);
        //改变窗口大小时，重置弹窗的宽高，防止超出可视区域（如F12调出debug的操作）
        $(window).on("resize",function(){
            layui.layer.full(index);
        })
    }
    $(".addNews_btn").click(function(){
        addNews();
    })

    //批量删除
    $(".delAll_btn").click(function(){
        var checkStatus = table.checkStatus('newsListTable'),
            data = checkStatus.data,
            newsId = [];
        if(data.length > 0) {
            for (var i in data) {
                newsId.push(data[i].newsId);
            }
            layer.confirm('确定删除选中的文章？', {icon: 3, title: '提示信息'}, function (index) {
                // $.get("删除文章接口",{
                //     newsId : newsId  //将需要删除的newsId作为参数传入
                // },function(data){
                tableIns.reload();
                layer.close(index);
                // })
            })
        }else{
            layer.msg("请选择需要删除的文章");
        }
    })

    //列表操作
    table.on('tool(newsList)', function(obj){
        var layEvent = obj.event,
            data = obj.data;

        if(obj.event === 'edit'){
            layer.prompt({
                title:obj.data.user+'输入新的密码',
                formType: 3
                ,value: ''
            }, function(value, index){
                var $ = layui.jquery;//这里很关键 加载jquery 否则不执行
                $.ajax({
                    type: 'post',
                    url: '/index/user/mima',
                    data: {
                        id: data.id,
                        ippass:value,
                    },
                    dataType: 'json',
                    success: function(data) {

                        if (data.status=="200"){
                            layer.msg('密码修改成功');
                        }else {
                            layer.msg(data.info);
                        }

                    }

                });
                return false;

            });
        }else if(obj.event === 'xufei'){
            layer.open({
                title:'账号：'+data.user+'---连接数：'+data.value,
                type: 2,
                area: ['340px', '340px'],
                content: '/index/user/xufei/id/'+data.id //这里content是一个URL，如果你不想让iframe出现滚动条，你还可以content: ['http://sentsin.com', 'no']
            });
        }else if(obj.event === 'tixian'){

            layer.open({
                title:'账号：'+data.user+'---连接数：'+data.value,
                type: 2,
                area: ['300px', '240px'],
                content: '/index/user/tixian/id/'+data.id //这里content是一个URL，如果你不想让iframe出现滚动条，你还可以content: ['http://sentsin.com', 'no']
            });
        }



    });

})