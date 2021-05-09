/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

// 双向链表
type DoubleList struct {
	Size uint
	Head *DoubleNode
	Tail *DoubleNode
}

// 节点数据
type NodeObject interface{}

// 双链表节点
type DoubleNode struct {
	Data NodeObject
	Prev *DoubleNode
	Next *DoubleNode
}

// tail->head
func (list *DoubleList) Append(data interface{}) {
	node := new(DoubleNode)
	node.Data = data
	if list.Size == 0 {
		list.Head = node
		list.Tail = node
		node.Data = data
	} else {
		list.Tail.Prev = node
		list.Tail.Prev.Next = list.Tail
		list.Tail = list.Tail.Prev
	}
	list.Size++
}

func (list *DoubleList) Delete(node *DoubleNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if node == list.Head {
		list.Head = node.Prev
	}
	if node == list.Tail {
		list.Tail = node.Next
	}
	list.Size--
}

func (list *DoubleList) DeleteHead(node *DoubleNode)  {
	oldHead :=list.Head
	if list.Size ==1{
		list.Head =nil
		list.Tail =nil
	}else{
		list.Head = oldHead.Prev
		list.Head.Next = nil
	}
}
