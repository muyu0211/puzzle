export default class Puzzle {
    constructor() {
        // 原始图像的宽高
        this.width = 0;
        this.height = 0;
        this.originImage = null
        this.shuffleImage = null

        // 存储打乱后的图片块编号
        this.arrImage = []
        this.cols = 0
        this.rows = 0

        // 左边拼图块
        this.puzzle = document.getElementById('left-son')
        // 右边拼图接收区
        this.puzzleDestionation = document.getElementById("right-son")

        // 重新开始按钮
        this.gameRestartButton = document.getElementById("gameRestartButton")
        // 查看原图按钮
        this.checkOriginButton = document.getElementById('checkOriginButton');
        // 模态弹窗
        this.modal = document.getElementById('imageModal');
        // 模态弹窗中的图片元素
        this.modalImage = document.getElementById('modalImage');
        // 模特窗口的关闭组件
        this.modelClose = document.getElementById('modelClose');
        // 查看拼图块
        this.checkShuffleButton = document.getElementById('checkShuffleButton');
        // 保存拼图块
        this.saveShuffleButtion = document.getElementById('saveShuffleButton')
        // 提交拼图按钮
        this.submitButton = document.getElementById('submitButton');

        // 烟花
        this.jsConfetti = new JSConfetti()

        this.init()
    }

    // 初始化烟花特效
    initConfetti(i = 1){
        for (i; i > 0; i--) {
            setTimeout(() => {
                this.jsConfetti.addConfetti()
            }, i * 500)
        }
    }

    async init() {
        // 初始化拼图（异步操作），先执行
        await this.initPuzzle();

        // 监听是否重新游戏
        this.addEventListenerIsReStartBtn();

        // 监听是否查看原图
        this.addEventListenerCheckOriginBtn();

        // 监听是否关闭模态窗口
        this.addEventListenerCloseModelBtn()

        // 监听是否查看拼图块
        this.addEventListenerCheckShuffleButton()

        // 监听提交拼图按钮
        this.addEventListenerSubmitPuzzleBtn()
    }

    // 初始化拼图
    async initPuzzle() {
        let loadingDiv = document.getElementById("loading")
        loadingDiv.style.display = "block"

        // 向后端请求原图数据
        try {
            const response = await fetch('/getImage');
            if (response.ok) {
                const data = await response.json();
                // 获取 Base64 编码的图片字符串
                this.originImage = data.image;
                this.arrImage = data.randomImages;
                this.cols = data.cols;
                this.rows = data.rows

                // 加载 Base64 编码的图片并获取其尺寸
                const img = new Image();
                // 将 Base64 编码的图片作为图片源
                img.src = this.originImage;
                img.onload = () => {
                    // 图片加载完成后，获取其宽高
                    this.width = img.width;
                    this.height = img.height;
                    console.log(`Width: ${this.width}, Height: ${this.height}`);

                    // 创建左边拖拽区
                    this.createPuzzle(this.puzzle)
                    // 创建右边接收区
                    this.createPuzzle(this.puzzleDestionation)

                    // 移除加载提示
                    loadingDiv.style.display = "none"
                };
            } else {
                console.error('Failed to fetch image data:', response.status);
            }
        } catch (error) {
            console.error('Error fetching image data:', error);
        }

        // 像后端请求拼图数据
        try {
            // 向后端发起请求获取 Base64 编码的图片数据
            const response = await fetch('/getShuffles');
            if (response.ok) {
                const data = await response.json();
                this.shuffleImage = data.image;

                const img = new Image();
                img.src = this.shuffleImage;
                img.onload = () => {
                    this.width = img.width;
                    this.height = img.height;
                    console.log(`Width: ${this.width}, Height: ${this.height}`);
                    loadingDiv.style.display = "none"
                };
            } else {
                console.error('Failed to fetch shuffle data:', response.status);
            }
        } catch (error) {
            console.error('Error fetching shuffle data:', error);
        }
    }

    // 监听查看原图按钮点击事件
    addEventListenerCheckOriginBtn() {
        this.checkOriginButton.addEventListener('click', (event) => {
            event.preventDefault();
            this.modalImage.src = this.originImage
            // 显示模态弹窗
            this.saveShuffleButtion.style.display = "none"
            this.modal.style.display = "block";
        });
    }

    // 监听关闭模特窗口
    addEventListenerCloseModelBtn() {
        this.modelClose.addEventListener('click', () => {
            this.modal.style.display = "none";
        });
    }

    // 监听是否点击重新游戏按钮
    addEventListenerIsReStartBtn(){
        // 获取元素
        const overlay = document.getElementById('confirmOverlay');
        const cancelBtn = document.getElementById('cancelBtn');
        const confirmBtn = document.getElementById('confirmBtn');

        // 点击重新开始按钮时，显示确认框
        this.gameRestartButton.addEventListener("click", () => {
            // window.location.reload()
            overlay.style.display = 'flex';  // 显示模态框
        })

        // 点击取消按钮时，隐藏确认框
        cancelBtn.addEventListener('click', function() {
            overlay.style.display = 'none';  // 隐藏模态框
        });

        // 点击确认按钮时，执行删除操作
        confirmBtn.addEventListener('click', function() {
            window.location.reload()
            overlay.style.display = 'none';  // 隐藏模态框
        });
    }

    // 监听是否查看拼图块
    addEventListenerCheckShuffleButton() {
        let loadingDiv = document.getElementById("loading")
        loadingDiv.style.display = "block"

        this.checkShuffleButton.addEventListener("click", async (event) => {
            event.preventDefault();
            this.modalImage.src = this.shuffleImage
            this.modal.style.display = "block";
            this.saveShuffleButtion.style.display = "block"
        })

        this.saveShuffleButtion.addEventListener("click", () => {
            // 解析 Base64 数据
            let blob = base64ToBlob(this.shuffleImage);

            // 创建下载链接
            let downloadLink = document.createElement('a');
            downloadLink.href = URL.createObjectURL(blob);
            downloadLink.download = 'image.png'; // 下载文件名为 image.png

            // 模拟点击下载链接
            downloadLink.click();

            // 释放 URL 对象
            URL.revokeObjectURL(downloadLink.href);

            // 将 Base64 数据解析为 Blob 对象
            function base64ToBlob(base64data) {
                let arr = base64data.split(',');
                let mime = arr[0].match(/:(.*?);/)[1];
                let bstr = atob(arr[1]);
                let n = bstr.length;
                let u8arr = new Uint8Array(n);

                while (n--) {
                    u8arr[n] = bstr.charCodeAt(n);
                }
                return new Blob([u8arr], { type: mime });
            }
        })
    }

    // 监听提交拼图按钮
    addEventListenerSubmitPuzzleBtn() {
        const overlay = document.getElementById("uncompletedOverlay")
        const confirmBtn = document.getElementById("uncompletedBtn")

        // 点击提交拼图按钮时，检测是否所有拼图块都被接收区接收，若是，则提交各拼图块id给后端，进行顺序判断；否则显示未完成提示框
        this.submitButton.addEventListener("click", async () => {
            let left = document.getElementById("left-son")
            let right = document.getElementById("right-son")
            let textElm = document.getElementById("uncompletedText")
            let loadingDiv = document.getElementById("loading")
            let allPiece = left.childNodes
            // 所有拼图块均被接收，则可以提交
            if (allPiece.length !== 0) {
                textElm.textContent = "拼图尚未完成"
                overlay.style.display = 'flex';  // 显示模态框
            } else {
                // TODO: 提交各拼图块id给后端，进行顺序判断
                allPiece = right.childNodes
                let pieceIds = []
                allPiece.forEach(pieceDiv => {
                    let piece = pieceDiv.childNodes[0]
                    // TODO: 提交各拼图块id给后端，进行顺序判断
                    if (piece !== undefined) {
                        pieceIds.push(piece.id)
                    }
                })

                loadingDiv.style.display = "block"
                // 向后端发起响应
                await fetch("/checkPuzzle", {
                    method: "POST",
                    body: JSON.stringify({
                        puzzleImages: pieceIds
                    })
                })
                    .then(response => response.json())
                    .then(data => {
                        // 移除加载提示
                        loadingDiv.style.display = "none"

                        if (data.data.code !== 200) {
                            textElm.textContent = "拼图尚未完成"
                            overlay.style.display = 'flex';
                        } else {
                            textElm.textContent = "恭喜你，完成拼图!"
                            overlay.style.display = 'flex';  // 显示模态框
                            this.initConfetti(10)
                        }

                    })
            }
        })

        // 点击确认按钮时，执行提交操作
        confirmBtn.addEventListener('click', function() {
            overlay.style.display = 'none';  // 隐藏模态框
        });
    }

    // 图像比例缩放
    scaleImage() {
        // 缩放比例 (调整为适合屏幕的大小，比如以最大宽度或高度为基准)
        let maxContainerWidth = 1000; // 设定容器最大宽度
        let maxContainerHeight = 580; // 设定容器最大高度
        let scale
        if (this.width > this.height || this.width > maxContainerWidth){
            scale = Math.min(maxContainerWidth / this.width, 1);
        }else {
            scale = Math.min(maxContainerHeight / this.height, 1);
        }
        return scale
    }

    // 根据传入的puzzle,创建拖拽元素
    createPuzzle(ele) {
        // 缩放比例 (调整为适合屏幕的大小，比如以最大宽度或高度为基准)
        let scale = this.scaleImage()
        // 缩放后的宽高
        this.width = this.width * scale;
        this.height = this.height * scale;

        // 设置偏移，使得拼图块居中显示
        let parentRight = document.getElementById("right-son");
        let offsetWidthRight = (parentRight.getBoundingClientRect().width - this.width) / 2;
        let offsetHeightRight = (parentRight.getBoundingClientRect().height - this.height) / 2;

        // 存储打乱后的div
        let arr = [];
        // 给每个div设置一个唯一id
        let id_num = 0;
        // 每个puzzle块的宽高
        let w = `${Math.floor(this.width / this.cols)}`
        let h = `${Math.floor(this.height / this.rows)}`
        // 根据用户选择的行列数
        for (let i = 0; i < this.rows; i++) {
            for (let j = 0; j < this.cols; j++) {
                let div = document.createElement("div");
                div.style.backgroundRepeat = "no-repeat";
                div.style.zIndex = '50'
                div.style.width = w + 'px'
                div.style.height = h + 'px'
                div.style.position = 'absolute'
                div.style.maxHeight = '100%'
                div.style.boxSizing = 'border-box'
                if (ele.id === 'left-son') {
                    // 映射得到一维数组下标，并计算获得打乱后的二维坐标
                    let idx = i * this.cols + j
                    let r =  Math.floor(this.arrImage[idx] / this.cols);
                    let c = this.arrImage[idx] % this.cols

                    div.style.border = `${1}px solid black`
                    div.style.width = w + 'px'
                    div.style.height = h + 'px'
                    // 设置div位置
                    div.style.top = (i * h) + 'px'
                    div.style.left = (j * w) + 'px'

                    // 设置div背景图片
                    div.style.background = `url(${this.originImage})`
                    div.style.backgroundSize = `${this.width}px, ${this.height}px`

                    // 设置拼图块的背景图（取整张图像的固定块），并打乱顺序
                    div.style.backgroundPosition = `-${(c * w)}px -${(r * h)}px`
                    div.draggable = true                    // 设置可拖动
                    div.id = this.arrImage[idx]             // 虽然位置是乱的，但是id是正常的
                    arr.push(div)
                    // console.log("r, c:", r, c, div.id)
                } else {
                    div.style.border = `${1}px dotted #f6f9fe`
                    div.style.top = (i * h + offsetHeightRight) + 'px'
                    div.style.left = (j * w + offsetWidthRight) + 'px'
                    div.style.backgroundColor = 'rgba(0, 0, 0, 0.2)'
                    div.style.backgroundSize = `${this.width}px ${this.height}px`
                    div.style.backgroundPosition = `-${(j * w)}px -${(i * h)}px`
                    div.id = "re" + id_num
                    arr.push(div)
                }
                id_num++;
            }
        }
        // 设置拖拽事件和放置事件
        this.setPuzzleEvent(ele, arr)
    }

    // 设置拼图块拖拽事件
    setPuzzleEvent(ele, arr) {
        // 左边拼图拖拽的时候，设置id
        if (ele.id === 'left-son') {
            for (let i = 0; i < arr.length; i++) {
                ele.appendChild(arr[i])
                // 给每个拼图块绑定拖拽事件
                arr[i].ondragstart = (e) => {
                    e.dataTransfer.setData("Id", e.target.id)
                }
            }
        } else {
            // 遍历接收区的每个div
            for (let i = 0; i < arr.length; i++) {
                // let that = this
                ele.appendChild(arr[i])
                arr[i].ondragover = (e) => {
                    e.preventDefault()
                }

                arr[i].ondrop = function (e) {
                    e.preventDefault();                         // 防止浏览器默认行为
                    // 当前拖拽的拼图块
                    // TODO BUG_1: 左边拼图块和右边接收区div的id一样，在某些情况下会出错
                    let pieceId = e.dataTransfer.getData("Id");
                    let pieceDiv = document.getElementById(pieceId);

                    // 判断当前接收区div是否以及存在子元素（放置了拼图）
                    if (this.hasChildNodes()){
                        // 如果已有拼图块，记录当前目标 `div` 的拼图块
                        let existingDiv = this.firstChild;

                        // 交换位置：将原有拼图块放回拖拽拼图块的位置
                        let parentDiv = pieceDiv.parentNode; // 获取拖拽拼图块的原始父节点
                        parentDiv.appendChild(existingDiv);              // 将原有拼图块放回拖拽拼图块的位置
                        if (parentDiv.id === "left-son") {
                            pieceDiv.style.width = parseInt(pieceDiv.style.width) * scaleLeft + "px"
                            pieceDiv.style.height = parseInt(pieceDiv.style.height) * scaleLeft + "px"
                        }
                    }

                    pieceDiv.style.top = "0";
                    pieceDiv.style.left = "0";

                    // 将拖拽方块，从左边放到右边
                    this.appendChild(pieceDiv);
                }
            }
        }
    }
}